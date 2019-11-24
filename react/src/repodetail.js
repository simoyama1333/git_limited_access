import React, { Component } from 'react';
import axios from 'axios'
import Code from 'react-code-prettify';
import doc from './doc.png';
import folder from './folder.png';
import LoadingOrError from './loading';
import NoMatch from './404'

const sortBoolDesc = (ary,key) =>{
  ary.sort((a,b) =>{
      if(a[key]<b[key]) return -1;
      if(a[key]>b[key]) return 1;
      return 0;
  });
  return ary;
}

export default class RepoDetail extends Component {
    constructor(props){
      super(props);
      this.repon =  this.props.match.params.reponame;
      this.clickLimited = false;
      this.state = {value: '', needauth: true, repotree: null,failed: false,loading: true};
      this.handleSubmit = this.handleSubmit.bind(this);
      this.handleChange = this.handleChange.bind(this);
      this.handleClick = this.handleClick.bind(this);
      //this.fetchFile = this.fetchFile.bind(this);
    }

    handleChange(event) {
      this.setState({value: event.target.value});
    }
    handleSubmit(event) {
      event.preventDefault();
      axios.post('/api/auth',{
        reponame: this.repon,
        password: this.state.value,
        token: null
      })
      .then(response => {
        let newtoken = response.data.newtoken;
        this.setState({needauth: false});
        window.localStorage.setItem(this.repon + "_token",newtoken)
     // catchでエラー時の挙動を定義する
      }).catch(err => {
        if(err.response.status == 401){
          return this.setState({failed: true});
        }
        this.setState({
          err :  true
        });
        console.log('err:', err);
      });
    }
    //click
    handleClick(event){
      if(this.clickLimited == true){
        return;
      }
      this.clickLimited  = true;

      let type = event.target.getAttribute("type");
      //dirの場合開閉
      if(type == "dir"){
        let children = event.target.children;
        if(children < 2){
          return;
        }
        try{
          if(children[1].style.display != 'block'){
            children[1].style.display = 'block';
          }else{
            children[1].style.display = 'none';
          };
        }catch{}
      }else{
        let path = event.target.getAttribute("path");
        this.requestFile(path);

      }
      setTimeout( ()=>{this.clickLimited  = false} ,100);
    }
    
    componentDidMount (){
      let token = localStorage.getItem(this.repon + "_token");
      axios.post('/api/auth',{
        reponame: this.repon,
        password: null,
        token: token
      })
      .then(response => {
        let authresult = response.data.authresult;
        let code = response.data.readme;
        let path = response.data.path;
        let tree = response.data.tree;
        if(authresult  == true){
          this.setState({needauth: false,
             code: code,
             path: path,
             tree: tree,
             token: token,
             loading: false
            });
        }
       // catchでエラー時の挙動を定義する
      }).catch(err => {
        if(err.response.status == 404){
          this.setState({
            notfound :  true,
            loading: false
          });
          return 
        }
        this.setState({
          err :  true
        });
        console.log('err:', err);
      });
    }
    
    requestFile(path){
      console.log(path);
      console.log(this.state.token);
      axios.post('/api/request',{
        reponame: this.repon,
        token: this.state.token,
        path: path
      })
      .then(response => {
        console.log(response.data);
        this.setState({code: response.data.code,
          path: path
        });
       // 失敗時は表示しない
      }).catch(err => {
        this.setState({
          err :  true
        });
        console.log('err:', err);
      });
    }
  
    repoAuth(){
      return (
        <div>
        <form onSubmit={this.handleSubmit}>
          <label>
            Password:
            <input type="password" value={this.state.value} onChange={this.handleChange} />
          </label>
          <input type="submit" value="Submit" />
        </form>
        {this.state.failed && <div><br></br>Password is incorrect.</div>}
        </div>
      );
    }
    codeView(){
      if(this.state.code == ''){
        return NoMatch('File');
      }
      //()つきurlをリンクにする
      var reg = new RegExp("\\(((https?|ftp)(:\/\/[-_.!~*\'()a-zA-Z0-9;\/?:\@&=+\$,%#]+))\\)","g");
      let codeURLReplaced = this.state.code.replace(reg,"<a href='$1' target='_blank'>$1</a>");
      return(
        <div className="code">
          <img src={doc} className="icon"></img>
          {this.state.path}
          <Code codeString={codeURLReplaced}/>
          </div>
      )
    }
    treeView(){
        //再帰でtree構造を作る
        const dig = (ary,i) => {
          const sorted = sortBoolDesc(ary,"TypeFile");
          let list = [];
          i += 1;
          let first = i == 1 ? "" : "dir"
          sorted.map((item) =>{
            let digdig = dig(item.Files);
            list.push(
              <span>
              { item.TypeFile && (
            <li onClick={this.handleClick} type="file" path={item.Path}>
            <img src={doc} className="icon"></img>
            {item.Name}
            </li>
          )}
          {!item.TypeFile && (
            <li onClick={this.handleClick} type="dir" path={item.Path}>
            <img src={folder} className="icon"></img>
            {item.Name}
            {digdig}
            </li>
          )}
              </span>
            );
          }) 
          return (<ul className={first}>{list}</ul>)
        }
        let list = dig(this.state.tree,0);
        return(<div>
          {list}
          {this.codeView()}
          </div>
          )
    }
    render(){
      if(this.state.loading){
        return LoadingOrError(true);
      }
      if(this.state.notfound){
        return NoMatch();
      }
      if(this.state.needauth){
        return(<div><h2>{this.repon}</h2>{this.repoAuth()}</div>);
      }else{
        return(<div><h2>{this.repon}</h2>{this.treeView()}</div>)
      }
      
    }
  }
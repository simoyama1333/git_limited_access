import React, { Component } from 'react';
import axios from 'axios'
import Code from 'react-code-prettify';
import doc from './doc.png';
import folder from './folder.png';
import LoadingOrError from './loading';

export default class RepoDetail extends Component {
    constructor(props){
      super(props);
      this.repon =  this.props.match.params.reponame;
      this.state = {value: '', needauth: true, repotree: null,failed: false};
      this.handleSubmit = this.handleSubmit.bind(this);
      this.handleChange = this.handleChange.bind(this);
    }
    handleChange(event) {
      this.setState({value: event.target.value});
    }
    handleSubmit(event) {
      event.preventDefault();
      /* axios post
      success, if newtoken != ""
      localStrage.setItem(this.repon + "_token",token);
      this.setState({needauth: false});
      else
      this.setState({failed: true});
      end
  
      */
    }
    appearDisappear(event){
      if(event.target.children.length < 2){
        return;
      }
      if(event.target.children[1].style.display == 'none'){
        event.target.children[1].style.display = 'block';
      }else{
        event.target.children[1].style.display = 'none';
      };
    }
    
    componentDidMount (){
      let token = localStorage.getItem(this.repon + "_token");
      if(token){
        console.log("token aruyo");
        /* token check ok
          this.setState({needauth: false});
          localStrage.setItem(this.repon + "_token",token);
        */
      }
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
      const codeString = `function map(f, a) {
        var result = [], // Create a new Array
            i;
        for (i = 0; i != a.length; i++)
          result[i] = f(a[i]);
        return result;
      };`
      return(
        <div className="code">
          <img src={doc} className="icon"></img>
          README.mb
          <Code codeString={codeString}/>
          </div>
      )
    }
    treeView(){
        return(<div><ul>
            <li onClick={this.appearDisappear}><img src={folder} className="icon"></img>親リスト1
            <ul className="dir">
              <li><img src={doc} className="icon"></img>子リスト1</li>
              <li><img src={folder} className="icon"></img>子リスト2
              <ul className="dir">
              <li>子リスト31</li>
              <li>子リスト32</li>
              <li>子リスト33</li>
            </ul>
          
              </li>
            </ul>
            </li>
            <li>親リスト2</li>
            <li>親リスト3
          
            </li>
          </ul>
          {this.codeView()}
          </div>
          )
    }
    render(){
      if(this.state.needauth){
        return(<div><h2>{this.repon}</h2>{this.repoAuth()}</div>);
      }else{
        return(<div><h2>{this.repon}</h2>{this.treeView()}</div>)
      }
      
    }
  }
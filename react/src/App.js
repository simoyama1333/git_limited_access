import lock from './lock.png';
import doc from './doc.png';
import folder from './folder.png';
import './App.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch} from 'react-router-dom';
import axios from 'axios'




const About = () => (
  <div><h2>About</h2></div>
)
const Topics = () => {
  return (
<div><h2>Topics</h2></div>
  )
}
const LoadingOrError = (errflag) => {
  if(errflag){
    return (<div>Loading...</div>)
  }else{
    return (<div>Error occurred</div>)
  }
}
const NoMatch = () =>{
  return (
    <div>
      <h2>404 Not Found</h2>
    </div>
  );
}
const testtest = () =>{
  console.log(this);
}

class RepoDetail extends Component {
  constructor(props){
    super(props);
    this.repon =  this.props.match.params.reponame;
    this.state = {value: '', needauth: true, repotree: null,failed: false};
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.test = {

    }
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

  repoauth(){
    return (
      <div>
      <form onSubmit={this.handleSubmit}>
      <h2>{this.repon}</h2>
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
  render(){
    if(this.state.needauth){
      return(<div>{this.repoauth()}
      <ul>
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
      
      </div>);
    }
    
  }
}
class Repos extends Component {
  constructor(props) {
    super(props);
    this.state = { repolist : null ,username : null, err : null}
  }
  componentDidMount (){
    axios.get('/api/repolist')
    .then(response => {
      console.log(response.data);
      this.setState({
        repolist :  response.data.repolist,
        username :  response.data.username
      });

   // catchでエラー時の挙動を定義する
    }).catch(err => {
      this.setState({
        err :  true
      });
      console.log('err:', err);
    });
  }

  render(){
    if(this.state.repolist == null){
      return LoadingOrError(false);
    }
    if(this.state.err != null){
      return LoadingOrError(true);
    }
    var list = [];
    {this.state.repolist.map((i) => {
      list.push(
          <li>
            {i.password_flag && <img src={lock} className="icon"></img>}
            <Link to={"/repodetail/" + i.name}>{i.name}</Link> 
            &nbsp;&nbsp;&nbsp;  
            {i.expire_flag && <span>閲覧期限:{i.expire}まで</span>}
          </li>
          );
    })}

    return (
      <div>
        <h3>{this.state.username}'s Repository</h3>
        {list}
      </div>
      )
    
  }
}
class GitHubLink extends Component {
  constructor(props) {
    super(props);
    this.state = { link : null }
  }
  componentDidMount (){
    axios.get('/api/username')
    .then(response => {
      console.log(response.data);
      this.setState({link : "https://github.com/" + response.data.username });

   // catchでエラー時の挙動を定義する
    }).catch(err => {
      console.log('err:', err);
    });
  }
  render(){
    if(this.state.link){
      return (<a href={this.state.link}>GitHub</a>)
    }else{
      return (<div></div>)
    }
    
  }
}

class App extends Component {
  render() {
  return (
    <Router>
        <div>
        <nav>
        <Link to="/">Home</Link>
        <GitHubLink />
      </nav>
      <br></br>
          <Switch>
            <Route exact path="/" component={Repos} />
            <Route exact path="/repo/:reponame" component={RepoDetail} />
            <Route exact path="/about" component={About} />
            <Route component={NoMatch}/>
          </Switch>
        </div>
      </Router>
  );
  }
}

export default App;

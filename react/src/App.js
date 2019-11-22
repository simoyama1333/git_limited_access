import lock from './lock.png';
import './App.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link} from 'react-router-dom';
import axios from 'axios'


const About = () => (
  <div><h2>About</h2></div>
)
const Topics = () => {
  return (
<div><h2>Topics</h2></div>
  )
}
class Repos extends Component {
  constructor(props) {
    super(props);
    this.state = { repolist : null ,username : null}
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
      console.log('err:', err);
    });
  }
  render(){
    if(this.state.repolist == null){
      return (<div>loading...</div>)
    }else{
      var list = [];
      {this.state.repolist.map((i) => {
        console.log(i);
        list.push(
          <li>
            {i.password_flag && <img src={lock} class="lock"></img>}
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
            <Route exact path="/" component={Repos} />
            <Route exact path="/about" component={About} />
            <Route exact path="/topics" component={Topics} />
        </div>
      </Router>
  );
  }
}

export default App;

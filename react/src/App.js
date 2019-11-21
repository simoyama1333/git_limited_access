import logo from './logo.svg';
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
    this.state = { link : "aaa" }
  }
  componentDidMount (){
    axios.get('https://jsonplaceholder.typicode.com/posts')
    .then(response => {
      console.log(response.data);
      this.setState({aja : "ue-i"});

   // catchでエラー時の挙動を定義する
    }).catch(err => {
      console.log('err:', err);
    });
  }
  render(){
    return (<div>{this.state.link}</div>)
  }
}
class GitHubLink extends Component {
  constructor(props) {
    super(props);
    this.state = { link : "" }
  }
  componentDidMount (){
    axios.get('https://localhost/')
    .then(response => {
      console.log(response.data);
      this.setState({aja : "ue-i"});

   // catchでエラー時の挙動を定義する
    }).catch(err => {
      console.log('err:', err);
    });
  }
  render(){
    return (<div>{this.state.link}</div>)
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

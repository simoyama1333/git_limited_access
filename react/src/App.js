import './App.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Link, Switch} from 'react-router-dom';
import axios from 'axios'
import Code from 'react-code-prettify';
import doc from './doc.png';
import folder from './folder.png';
import RepoDetail from './repodetail';
import Repos from './repo';
import NoMatch from './404'

const About = () => (
  <div><h2>About</h2></div>
)
const Topics = () => {
  return (
<div><h2>Topics</h2></div>
  )
}


const testtest = () =>{
  console.log(this);
}

class GitHubLink extends Component {
  constructor(props) {
    super(props);
    this.state = { link : null }
  }
  componentDidMount (){
    axios.get('/api/username')
    .then(response => {
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

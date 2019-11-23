import React, { Component } from 'react';
import axios from 'axios'
import { BrowserRouter as Link} from 'react-router-dom';
import lock from './lock.png';
import LoadingOrError from './loading';

export default class Repos extends Component {
    constructor(props) {
      super(props);
      this.state = { repolist : null ,username : null, err : null}
    }
    componentDidMount (){
      axios.get('/api/repolist')
      .then(response => {
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
              <a href={"/repo/" + i.name}>{i.name}</a> 
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
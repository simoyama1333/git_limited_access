import React from 'react';

export default (errflag) => {
    if(errflag){
      return (<div>Loading...</div>)
    }else{
      return (<div>Error occurred</div>)
    }
  }
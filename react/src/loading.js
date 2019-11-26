import React from 'react';

export default (loading) => {
    if(loading){
      return (<div>Loading...</div>)
    }else{
      return (<div>Error occurred</div>)
    }
  }
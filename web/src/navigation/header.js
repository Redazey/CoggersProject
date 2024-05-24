import React, { useState } from "react";
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import './header.css'; 

const Header = () => {
  return (
    <Router>
    <div className="Header">
        <div className="logo">
            <Link className="logo__img" to="#"><img src={require('../imgs/logo.png')}/></Link>
            <div className="logo__text"></div>
        </div>
        
        <div className="MainButtons">
            <Link to="#">Играть</Link>
            <Link to="#">Войти</Link>
        </div>
    </div>
    </Router>
    );
};

export default Header;

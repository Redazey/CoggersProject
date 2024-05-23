import React, { useState } from "react";
import './header.css'; 

const Header = () => {
  return (
    <div className="Header">
        <div className="logo">
            <a className="logo__img" href="#"><img src={require('../imgs/logo.png')}/></a>
            <div className="logo__text"></div>
        </div>
        
        <div className="MainButtons">
            <a href="#">Играть</a>
            <a href="#">Войти</a>
        </div>
    </div>
    );
};

export default Header;

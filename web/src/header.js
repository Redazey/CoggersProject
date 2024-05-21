import React, { useState } from "react";
import './header.css'; 

const Header = () => {
  return (
    <div className="Header">
        <div className="horizontalNav">
            <a className="horizontalNav__logo" href="#"><img src={require('./imgs/logo.png')}/></a>
            <div className="horizontalNav__MainButtons">
                <a href="#">Играть</a>
                <a href="#">Войти</a>
            </div>
        </div>

        <div className="verticalNav">
            <ul class="vertical-list">
                <li><a href="#">Личный кабинет</a></li>
                <li><a href="#">Помощь</a></li>
                <li><a href="#">Магазин</a></li>
            </ul>
        </div>
    </div>
    );
};

export default Header;

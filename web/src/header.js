import React, { useState } from "react";
import './header.css'; 

const Header = () => {
  return (
    <div className="Header">
      <nav>
        <a href="#" className="navHref">Личный кабинет</a>
        <a href="#" className="navHref">Помощь</a>
        <a href="#" className="navStyledHref">Играть</a>
        <a href="#" className="navHref">Магазин</a>
        <a href="#" className="navStyledHref">Войти</a>
      </nav>
    </div>
  );
};

export default Header;
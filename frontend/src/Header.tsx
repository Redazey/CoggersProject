import React from "react";
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import logoImg from './imgs/logo.png';
import './header.css'; 

const Header: React.FC = () => {
  return (
    <Router>
      <div className="Header">
          <div className="logo">
              <Link className="logo__img" to="#"><img src={logoImg} alt="Logo"/></Link>
              <div className="logo__text"><h3>Coggers Project</h3></div>
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

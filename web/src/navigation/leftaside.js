import React, { useState } from "react";
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';
import './leftaside.css'; 

const Leftaside = () => {
  return (
    <Router>
    <aside>
        <ul class="vertical-list">
            <li><img src={require('../imgs/cab_logo.png')} alt="cab_img"/><Link to="#">Личный кабинет</Link></li>
            <li><img src={require('../imgs/shop_logo.png')} alt="shop_img"/><Link to="#">Магазин</Link></li>
            <li><img src={require('../imgs/help_logo.png')} alt="help_img"/><Link to="#">Помощь</Link></li>
        </ul>
    </aside>
    </Router>
    );
};

export default Leftaside;

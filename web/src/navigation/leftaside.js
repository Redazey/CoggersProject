import React, { useState } from "react";
import './leftaside.css'; 

const Leftaside = () => {
  return (
    <aside>
        <ul class="vertical-list">
            <li><img src={require('../imgs/cab_logo.png')} alt="cab_img"/><a href="#">Личный кабинет</a></li>
            <li><img src={require('../imgs/help_logo.png')} alt="help_img"/><a href="#">Помощь</a></li>
            <li><img src={require('../imgs/shop_logo.png')} alt="shop_img"/><a href="#">Магазин</a></li>
        </ul>
    </aside>
    );
};

export default Leftaside;

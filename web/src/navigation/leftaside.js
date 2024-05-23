import React, { useState } from "react";
import './leftaside.css'; 

const Leftaside = () => {
  return (
    <div className="sidebar">
        <ul class="vertical-list">
            <li><a href="#">Личный кабинет</a></li>
            <li><a href="#">Помощь</a></li>
            <li><a href="#">Магазин</a></li>
        </ul>
    </div>
    );
};

export default Leftaside;

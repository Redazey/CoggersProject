import React, { useState } from "react";
import './MainPage.css'; 
import { BrowserRouter as Router, Link } from 'react-router-dom';
import Header from "../../navigation/header";
import Leftaside from "../../navigation/leftaside";

const LoginPage = () => {
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");

    const handleLogin = () => {
        // Здесь можно добавить логику для обработки регистрации пользователя
        console.log("Вход пользователя:", { login, password });
    };

    const scrollImage = () => {
    window.addEventListener('scroll', function() {
        var img = document.getElementById('scrolling-image');
        var scrollHeight = window.scrollY / 3;
        img.style.transform = `rotate(${scrollHeight}deg)`;
    });
    }

    scrollImage();

    /*
    обязательно сделать news-item как функцию, что бы она подтягивала новости из бд
    */
    return (
    <div className="Content">
        <Header />
        <main>
            <div className="top-content">
                <img id="scrolling-image" src={require('../../imgs/content.png')}/>
                <scrollImage />
            </div>

            <div className="mid-content">
                <div className="news-container">
                    <div className="news-item">
                        
                        <div className="news-text">
                            <h3>Мы мобилизируем игроков!</h3>
                            <p>
                                В связи с грядущей второй волной мобилизации мы призываем всех игроков нашего сервера учавствовать в СВО!
                            </p>    
                        </div>
                        
                        <img src={require('../../imgs/new1.png')}/>
                    </div>
                </div>  
            </div>
        </main>
    </div>
    );
};

export default LoginPage;


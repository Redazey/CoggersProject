import React, { useState } from "react";
import './MainPage.css'; 
import Header from "../../navigation/header";
import Leftaside from "../../navigation/leftaside";

const LoginPage = () => {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    // Здесь можно добавить логику для обработки регистрации пользователя
    console.log("Вход пользователя:", { login, password });
  };

  return (
    <div className="Content">
        <Header />
        <Leftaside />
        <div className="main">
            <p>

            </p>
            <img />
            
        </div>
    </div>
  );
};

export default LoginPage;


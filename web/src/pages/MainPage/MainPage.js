import React, { useState } from "react";
import './MainPage.css'; 
import Header from "../../navigation/header";
import leftaside from "../../navigation/leftaside";

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
        <leftaside />
    </div>
  );
};

export default LoginPage;


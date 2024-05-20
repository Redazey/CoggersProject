import React, { useState } from "react";
import './login.css'; 
import Header from "../../header";

const LoginPage = () => {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    // Здесь можно добавить логику для обработки регистрации пользователя
    console.log("Вход пользователя:", { login, password });
  };

  return (
    <Header />
    
  );
};

export default LoginPage;


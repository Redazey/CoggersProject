import React, { useState } from "react";

const LoginPage = () => {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    // Здесь можно добавить логику для обработки регистрации пользователя
    console.log("Вход пользователя:", { login, password });
  };

  return (
    <div>
      <h1>Вход</h1>
      <input type="text" placeholder="Логин" value={login} onChange={(e) => setLogin(e.target.value)} />
      <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleLogin}>Войти</button>
    </div>
  );
};

export default LoginPage;


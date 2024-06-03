import React, { useState } from "react";

const RegistrationPage = () => {
  const [login, setLogin] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleRegistration = () => {
    // Здесь можно добавить логику для обработки регистрации пользователя
    console.log("Регистрация пользователя:", { login, email, password });
  };

  return (
    <div>
      <h1>Регистрация</h1>
      <input type="text" placeholder="Логин" value={login} onChange={(e) => setLogin(e.target.value)} />
      <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button onClick={handleRegistration}>Зарегистрироваться</button>
    </div>
  );
};

export default RegistrationPage;


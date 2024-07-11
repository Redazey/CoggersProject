import React from "react";
import { Link } from 'react-router-dom';

const Header: React.FC = () => {
  return (  
        <>
            <Link to="#">Играть</Link>
            <Link to="#">Войти</Link>
        </>
    )
};

export default Header;

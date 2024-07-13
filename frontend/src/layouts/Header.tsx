import React from "react";
import { Link } from 'react-router-dom';
import logoImg from '../assets/images/logo.png'
import Navigation from "../components/Navigation"

const Footer: React.FC = () => {
    return (
        <>
            <div className="header__logo">
                <Link to="#"><img src={logoImg} alt="logo"/></Link>
                <h3>Coggers Project</h3>
            </div>

            <nav>
                <Navigation />
            </nav>
        </>
    );
};

export default Footer;

import React from "react";
import { BrowserRouter as Router, Link } from 'react-router-dom';
import './footer.css'; 

const Footer: React.FC = () => {
    return (
        <Router>
            <div className="Footer">
                <div className="Contacts">
                    {/* Здесь вы можете добавить контактную информацию */}
                </div>
            </div>
        </Router>
    );
};

export default Footer;

import React, { useState, useEffect } from "react";
import './App.css';
import Header from "./Header";
import Footer from "./Footer";
import ScrollRotate from "./scripts/scroll-rotate"
import contentImg from './imgs/content.png';
import NewsItem from "./news";

const App: React.FC = () => {
    useEffect(() => {
        ScrollRotate();
    }, []);

    return (
        <div className="Content">
            <Header />
            <main>
                <div className="top-content">
                    <img id="scrolling-image" src={contentImg} alt="Scrolling Image" />
                </div>

                <div className="mid-content">
                    <div className="news-container">
                        <NewsItem />
                    </div>
                </div>
            </main>
            <Footer />
        </div>
    );
};

export default App;

import React, { useState, useEffect } from "react";
import '../assets/styles/pages/Home/Home.min.css';

const Home: React.FC = () => {
    const [news, setNews] = useState<any[]>([]);

    return (
        <>
            {news.map((item, index) => (
                <div className="news__item" key={index}>
                    <div className="news__text">
                        <h3>{item.title}</h3>
                        <p>{item.description}</p>
                    </div>
                    <img src={item.image} alt={item.title} />
                </div>
            ))}
        </>
    );
};

export default Home;


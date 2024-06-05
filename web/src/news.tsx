import React, { useState, useEffect } from "react";
import './App.css';

const NewsItem: React.FC = () => {
    const [news, setNews] = useState<any[]>([]);

    useEffect(() => {
        fetch('/')
            .then(response => response.json())
            .then(data => setNews(data))
            .catch(error => console.error('Error fetching news:', error));
    }, []);

    return (
        <div>
            {news.map((item, index) => (
                <div className="news-item" key={index}>
                    <div className="news-text">
                        <h3>{item.title}</h3>
                        <p>{item.description}</p>
                    </div>
                    <img src={item.image} alt={item.title} />
                </div>
            ))}
        </div>
    );
};

export default NewsItem;


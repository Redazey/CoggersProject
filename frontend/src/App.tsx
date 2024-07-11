import React, { useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import ScrollRotate from "./scripts/scroll-rotate"
import logoImg from './assets/images/logo.png'
import contentImg from './assets/images/content.png'
import Navigation from "./components/Navigation"
import Footer from "./components/Footer"
import ServerItem from "./components/ServerItem"
import Home from "./pages/Home"
import NewsItem from "./pages/NewsItem"

const App: React.FC = () => {

    useEffect(() => {
        ScrollRotate();
    }, []);

    return (
        <>
            <Router>
                <div className="container">
                    <header>
                        <div className="header__logo">
                            <Link to="#"><img src={logoImg} alt="logo"/></Link>
                            <h3>Coggers Project</h3>
                        </div>

                        <nav>
                            <Navigation />
                        </nav>
                    </header>

                    <main>
                        <div className="main__top">
                            <img id="scrolling-image" src={contentImg} alt="Scrolling Image" />
                        </div>

                        <div className="main__middle">
                            <ServerItem />
                            <Routes>
                                {/* Панель с серверами */}
                                <Route path="/server/id/:id" element={<ServerItem />} />
                                {/* Мэйн */}
                                <Route path="/" element={<Home />} />
                                <Route path="/news/id/:id" element={<NewsItem />} />
                                {/* Админ-панель */}
                            </Routes>
                        </div>
                    </main>
                    
                    <footer>
                        <Footer />
                    </footer>
                </div>
            </Router>
        </>
    );
};

export default App;

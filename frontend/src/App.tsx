import React, { useEffect, lazy } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import ScrollRotate from "./scripts/scroll-rotate"
import contentImg from './assets/images/content.png'
import Aside from "./layouts/Aside"
import Header from "./layouts/Header"
import Footer from "./layouts/Footer"

const Home = lazy(() => import("./pages/Home"));
const ServerItem = lazy(() => import("./pages/ServerItem"));
const NewsItem = lazy(() => import("./pages/NewsItem"));

const App: React.FC = () => {

    useEffect(() => {
        ScrollRotate();
    }, []);

    return (
        <>
            <Router>
                <div className="container">
                    <header>
                        <Header />
                    </header>

                    <main>
                        <div className="main__top">
                            <img id="scrolling-image" src={contentImg} alt="Scrolling Image" />
                        </div>

                        <div className="main__middle">
                            <aside>
                                <Aside />
                            </aside>
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

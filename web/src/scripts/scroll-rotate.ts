const ScrollRotate = (): void => {
    window.addEventListener('scroll', function() {
        const img: HTMLElement | null = document.getElementById('scrolling-image');
        const scrollHeight: number = window.scrollY / 3;
        if (img) {
            img.style.transform = `rotate(${scrollHeight}deg)`;
        }
    });
};

export default ScrollRotate;
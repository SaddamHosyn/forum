// Parallax Effect for Floating Posts
document.addEventListener('scroll', () => {
    const postsContainer = document.getElementById('posts');
    const scrollPosition = window.scrollY;

    // Apply the same transformation to all posts
    const speed = 0.3; // Adjust speed as needed
    const offset = scrollPosition * speed;

    // Limit the maximum and minimum offset to prevent posts from moving too far
    const maxOffset = 100; // Adjust as needed
    const minOffset = -50; // Adjust as needed
    const clampedOffset = Math.min(Math.max(offset, minOffset), maxOffset);

    postsContainer.style.transform = `translateY(${clampedOffset}px)`;
});
// document.addEventListener('scroll', () => {
//     requestAnimationFrame(() => {
//         const posts = document.querySelectorAll('.post-item');
//         const scrollPosition = window.scrollY; // Use window.scrollY instead of window.pageYOffset

//         posts.forEach((post, index) => {
//             // Adjust the vertical position of each post based on scroll
//             const speed = 0.3 + (index * 0.1); // Adjust speed for each post
//             post.style.transform = `translateY(${scrollPosition * speed}px)`; //Move downward
//         });
//     });
// });
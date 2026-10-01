// DVD Screensaver Emoji Box
document.addEventListener("DOMContentLoaded", () => {
    const stage = document.getElementById("emoji-stage");
    if (!stage) {
        return;
    }

    const emojiElements = stage.querySelectorAll(".dvd-emoji");

    // If there are no emojis, display the in-stage empty overlay
    if (emojiElements.length === 0) {
        const emptyOverlay = document.createElement("div");
        emptyOverlay.className = "stage-empty-overlay";
        emptyOverlay.innerHTML = `
            <span class="empty-icon">🚲</span>
            <p>No reactions yet.<br>Be the first to react!</p>
        `;
        stage.appendChild(emptyOverlay);
        return;
    }

    let stageWidth = stage.clientWidth;
    let stageHeight = stage.clientHeight;
    let animationFrameId = null;
    let isRunning = false;

    // Glow classes defined in emoji-box.css
    const glowClasses = [
        "bounce-glow-0",
        "bounce-glow-1",
        "bounce-glow-2",
        "bounce-glow-3"
    ];

    // Helper: random float in [min, max]
    function randomRange(min, max) {
        return Math.random() * (max - min) + min;
    }

    // Initialize state for each emoji
    const emojis = Array.from(emojiElements).map((el, index) => {
        const width = el.offsetWidth || 44;
        const height = el.offsetHeight || 44;

        // Spread initial positions across the stage
        const maxX = Math.max(0, stageWidth - width);
        const maxY = Math.max(0, stageHeight - height);
        const x = randomRange(0, maxX);
        const y = randomRange(0, maxY);

        // Calculate random velocity and angle (avoiding exact 45 degrees)
        const speed = randomRange(.25, 1);
        const angle = randomRange(0.4, 1.1); // approx 23 deg to 63 deg
        const dirX = Math.random() < 0.5 ? -1 : 1;
        const dirY = Math.random() < 0.5 ? -1 : 1;

        const vx = dirX * Math.cos(angle) * speed;
        const vy = dirY * Math.sin(angle) * speed;

        return {
            el,
            x,
            y,
            vx,
            vy,
            width,
            height,
            glowIndex: index % glowClasses.length
        };
    });

    function triggerBounce(item) {
        // Remove previous glow class
        glowClasses.forEach(cls => item.el.classList.remove(cls));

        // Advance to next glow color and apply
        item.glowIndex = (item.glowIndex + 1) % glowClasses.length;
        const newClass = glowClasses[item.glowIndex];
        item.el.classList.add(newClass);

        // Subtle squash/pop on bounce
        item.el.style.transform = `translate3d(${item.x}px, ${item.y}px, 0) scale(1.15)`;
        setTimeout(() => {
            if (item.el) {
                item.el.style.transform = `translate3d(${item.x}px, ${item.y}px, 0) scale(1)`;
            }
        }, 120);
    }

    // Animation physics loop
    function updatePhysics() {
        emojis.forEach(item => {
            item.x += item.vx;
            item.y += item.vy;

            let bounced = false;

            // Horizontal bounds check
            if (item.x <= 0) {
                item.x = 0;
                item.vx = Math.abs(item.vx);
                bounced = true;
            } else if (item.x + item.width >= stageWidth) {
                item.x = stageWidth - item.width;
                item.vx = -Math.abs(item.vx);
                bounced = true;
            }

            // Vertical bounds check
            if (item.y <= 0) {
                item.y = 0;
                item.vy = Math.abs(item.vy);
                bounced = true;
            } else if (item.y + item.height >= stageHeight) {
                item.y = stageHeight - item.height;
                item.vy = -Math.abs(item.vy);
                bounced = true;
            }

            if (bounced) {
                triggerBounce(item);
            } else {
                item.el.style.transform = `translate3d(${item.x}px, ${item.y}px, 0)`;
            }
        });

        if (isRunning) {
            animationFrameId = requestAnimationFrame(updatePhysics);
        }
    }

    function startAnimation() {
        if (!isRunning) {
            isRunning = true;
            animationFrameId = requestAnimationFrame(updatePhysics);
        }
    }

    function stopAnimation() {
        isRunning = false;
        if (animationFrameId) {
            cancelAnimationFrame(animationFrameId);
            animationFrameId = null;
        }
    }

    // Observe stage resizing to keep emojis inside bounds
    if (window.ResizeObserver) {
        const resizeObserver = new ResizeObserver(entries => {
            for (let entry of entries) {
                stageWidth = entry.contentRect.width;
                stageHeight = entry.contentRect.height;

                emojis.forEach(item => {
                    item.x = Math.min(Math.max(0, item.x), Math.max(0, stageWidth - item.width));
                    item.y = Math.min(Math.max(0, item.y), Math.max(0, stageHeight - item.height));
                });
            }
        });
        resizeObserver.observe(stage);
    }

    // Pause animation when offscreen to conserve CPU/battery
    if (window.IntersectionObserver) {
        const observer = new IntersectionObserver(entries => {
            entries.forEach(entry => {
                if (entry.isIntersecting && !document.hidden) {
                    startAnimation();
                } else {
                    stopAnimation();
                }
            });
        }, { threshold: 0.1 });
        observer.observe(stage);
    } else {
        startAnimation();
    }

    // Handle tab visibility changes
    document.addEventListener("visibilitychange", () => {
        if (document.hidden) {
            stopAnimation();
        } else {
            startAnimation();
        }
    });
});

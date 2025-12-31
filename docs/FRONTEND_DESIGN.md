### Part 1: Frontend Design Document

This document outlines the visual language and structure of the web application.

#### 1. Visual Identity
*   **Theme Name:** "Dark Mode Gamification" & "Light Mode Focus"
*   **Concept (Dark Mode):** A sleek, high-contrast interface that makes progress (green bars, checks) pop against a dark background. Feels like a mix between a social network and a gaming dashboard.
*   **Concept (Light Mode):** A clean, productivity-focused interface. Uses plenty of white space and subtle drop shadows to separate content, prioritizing readability and a "daytime" workflow feel.
*   **Typography:**
    *   **Headings:** `Poppins` (Bold, Modern, Geometric).
    *   **Body:** `Inter` (Clean, Highly Readable).
*   **Color Palette (Light Mode Variables):**
    *   `--bg-light`: #F3F4F6 (Cool Light Gray - Main Background)
    *   `--bg-card`: #FFFFFF (Pure White - Card Backgrounds)
    *   `--primary`: #4F46E5 (Rich Indigo - Branding/Links)
    *   `--accent`: #059669 (Deep Emerald - Completed/Success state)
    *   `--text-main`: #1F2937 (Dark Charcoal - Main Text)
    *   `--text-muted`: #6B7280 (Medium Gray - Secondary text)
    *   `--border`: #E5E7EB (Subtle Light Gray)

#### 2. Layout Structure (The "Holy Grail" Dashboard)
The page is divided into 3 distinct columns:
*   **Left Sidebar (Navigation):** Fixed position. Contains the Logo, Main Menu (Home, Lists, Challenges), and User Mini-Profile.
*   **Center Column (The Feed):** Scrollable. This is the social aspect. Contains status updates ("I just finished Chapter 5"), reviews, and milestone achievements.
*   **Right Column (Widgets):** Sticky. Contains "Trending Hobbies," "Friends' Activity," and "Current Streaks."

#### 3. Key UI Components
*   **Completion Card:** The core unit. Displays an item (e.g., Book Cover) + Title + Progress Bar + "Mark as Done" Button.
*   **The "Check" Badge:** A small circular icon that turns from gray to green when a task is done.
*   **Activity Feed Post:** A social post showing `[User Avatar] [Action] [Item Cover] [Caption]`.

---

### Part 2: The HTML Template (Light Mode)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Completionist | Track Everything</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Poppins:wght@600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    
    <style>
        :root {
            --bg-light: #F3F4F6;
            --bg-card: #FFFFFF;
            --primary: #4F46E5;
            --accent: #059669;
            --text-main: #1F2937;
            --text-muted: #6B7280;
            --border: #E5E7EB;
            --radius: 12px;
            --nav-width: 260px;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg-light);
            color: var(--text-main);
            line-height: 1.6;
        }

        a { text-decoration: none; color: inherit; transition: 0.2s; }
        ul { list-style: none; }

        .app-container {
            display: grid;
            grid-template-columns: var(--nav-width) 1fr 320px;
            gap: 24px;
            max-width: 1400px;
            margin: 0 auto;
            min-height: 100vh;
        }

        .sidebar {
            position: sticky;
            top: 0;
            height: 100vh;
            padding: 24px;
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
            background-color: var(--bg-light);
        }

        .logo {
            font-family: 'Poppins', sans-serif;
            font-size: 1.5rem;
            font-weight: 700;
            color: var(--text-main);
            margin-bottom: 40px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .logo i { color: var(--accent); }

        .nav-menu { display: flex; flex-direction: column; gap: 8px; }
        
        .nav-item {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 12px 16px;
            border-radius: var(--radius);
            color: var(--text-muted);
            font-weight: 500;
        }

        .nav-item:hover, .nav-item.active {
            background-color: #EEF2FF;
            color: var(--primary);
        }

        .nav-item i { width: 20px; text-align: center; }

        .user-mini-profile {
            margin-top: auto;
            display: flex;
            align-items: center;
            gap: 12px;
            padding-top: 20px;
            border-top: 1px solid var(--border);
        }

        .avatar {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            object-fit: cover;
            border: 2px solid var(--primary);
        }

        .user-info h4 { font-size: 0.9rem; color: var(--text-main); }
        .user-info span { font-size: 0.8rem; color: var(--accent); }

        .main-feed {
            padding: 24px 0;
        }

        .feed-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 24px;
        }

        .feed-title { font-family: 'Poppins', sans-serif; font-size: 1.25rem; }
        
        .create-post-btn {
            background-color: var(--primary);
            color: white;
            padding: 10px 20px;
            border-radius: 50px;
            font-weight: 600;
            font-size: 0.9rem;
            border: none;
            cursor: pointer;
            box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.3);
        }
        .create-post-btn:hover { opacity: 0.9; }

        .post-card {
            background-color: var(--bg-card);
            border-radius: var(--radius);
            padding: 20px;
            margin-bottom: 20px;
            border: 1px solid var(--border);
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
        }

        .post-header {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 16px;
        }

        .post-meta h4 { font-size: 0.95rem; }
        .post-meta span { font-size: 0.8rem; color: var(--text-muted); }

        .post-content { margin-bottom: 16px; color: var(--text-main); }

        .completion-item {
            display: flex;
            gap: 16px;
            background: #F9FAFB;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 16px;
            border-left: 4px solid var(--accent);
        }

        .item-img {
            width: 60px;
            height: 60px;
            border-radius: 6px;
            object-fit: cover;
        }

        .item-details h5 { font-size: 1rem; margin-bottom: 4px; }
        .tag {
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            padding: 2px 6px;
            border-radius: 4px;
            background: #E5E7EB;
            color: var(--text-muted);
        }
        .tag.anime { color: #DB2777; background: #FCE7F3; }
        .tag.book { color: #2563EB; background: #DBEAFE; }

        .completion-badge {
            display: flex;
            align-items: center;
            gap: 6px;
            color: var(--accent);
            font-weight: 600;
            font-size: 0.9rem;
            margin-top: 4px;
        }

        .post-actions {
            display: flex;
            gap: 20px;
            border-top: 1px solid var(--border);
            padding-top: 12px;
            color: var(--text-muted);
            font-size: 0.9rem;
        }

        .action-btn { display: flex; align-items: center; gap: 6px; cursor: pointer; }
        .action-btn:hover { color: var(--text-main); }

        .widgets {
            padding: 24px 0;
        }

        .widget-card {
            background-color: var(--bg-card);
            border-radius: var(--radius);
            padding: 20px;
            margin-bottom: 24px;
            border: 1px solid var(--border);
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
        }

        .widget-title {
            font-family: 'Poppins', sans-serif;
            font-size: 1rem;
            margin-bottom: 16px;
            color: var(--text-main);
        }

        .trending-list li {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 12px;
        }

        .trending-rank {
            font-weight: 700;
            color: var(--text-muted);
            width: 20px;
        }

        .trending-info { font-size: 0.9rem; }
        .trending-info span { display: block; font-size: 0.75rem; color: var(--text-muted); }

        @media (max-width: 1024px) {
            .app-container {
                grid-template-columns: 80px 1fr;
            }
            .sidebar { align-items: center; }
            .nav-item span, .user-info, .logo span { display: none; }
            .logo { margin-bottom: 20px; }
            .user-mini-profile { border-top: none; padding-top: 0; }
            .widgets { display: none; }
        }

        @media (max-width: 600px) {
            .app-container {
                display: block;
            }
            .sidebar {
                position: fixed;
                bottom: 0;
                top: auto;
                width: 100%;
                height: 60px;
                flex-direction: row;
                justify-content: space-around;
                padding: 0;
                background: var(--bg-card);
                border-top: 1px solid var(--border);
                z-index: 100;
                box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
            }
            .logo, .user-mini-profile { display: none; }
            .nav-menu { flex-direction: row; width: 100%; justify-content: space-around; }
            .nav-item { padding: 10px; font-size: 1.2rem; }
            .nav-item:hover { background-color: transparent; }
            .main-feed { padding: 16px; margin-bottom: 60px; }
        }
    </style>
</head>
<body>

    <div class="app-container">
        
        <aside class="sidebar">
            <div class="logo">
                <i class="fa-solid fa-circle-check"></i>
                <span>Completionist</span>
            </div>

            <nav class="nav-menu">
                <a href="#" class="nav-item active">
                    <i class="fa-solid fa-house"></i>
                    <span>Home Feed</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-list-check"></i>
                    <span>My Lists</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-compass"></i>
                    <span>Explore</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-trophy"></i>
                    <span>Challenges</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-book-open"></i>
                    <span>Library</span>
                </a>
            </nav>

            <div class="user-mini-profile">
                <img src="https://picsum.photos/seed/me/200/200" alt="User Avatar" class="avatar">
                <div class="user-info">
                    <h4>Alex Doe</h4>
                    <span>Lvl 42 • 500 Items</span>
                </div>
            </div>
        </aside>

        <main class="main-feed">
            
            <header class="feed-header">
                <h2 class="feed-title">Social Feed</h2>
                <button class="create-post-btn"><i class="fa-solid fa-plus"></i> Log Progress</button>
            </header>

            <article class="post-card">
                <div class="post-header">
                    <img src="https://picsum.photos/seed/sarah/200/200" alt="Sarah" class="avatar">
                    <div class="post-meta">
                        <h4>Sarah Jenkins</h4>
                        <span>2 hours ago</span>
                    </div>
                </div>
                
                <div class="post-content">
                    Finally finished this masterpiece! The animation in the last episode was insane. 10/10.
                </div>

                <div class="completion-item">
                    <img src="https://picsum.photos/seed/anime1/200/200" alt="Anime Cover" class="item-img">
                    <div class="item-details">
                        <span class="tag anime">Anime</span>
                        <h5>Neon Cyber Samurai</h5>
                        <div class="completion-badge">
                            <i class="fa-solid fa-circle-check"></i> Completed 12/12 Episodes
                        </div>
                    </div>
                </div>

                <div class="post-actions">
                    <div class="action-btn"><i class="fa-regular fa-heart"></i> 124</div>
                    <div class="action-btn"><i class="fa-regular fa-comment"></i> 18</div>
                    <div class="action-btn"><i class="fa-solid fa-share"></i> Share</div>
                </div>
            </article>

            <article class="post-card">
                <div class="post-header">
                    <img src="https://picsum.photos/seed/mike/200/200" alt="Mike" class="avatar">
                    <div class="post-meta">
                        <h4>Mike Ross</h4>
                        <span>5 hours ago</span>
                    </div>
                </div>
                
                <div class="post-content">
                    Halfway through "Dune". It's getting intense.
                </div>

                <div class="completion-item" style="border-left: 4px solid var(--primary);">
                    <img src="https://picsum.photos/seed/book1/200/200" alt="Book Cover" class="item-img">
                    <div class="item-details">
                        <span class="tag book">Book</span>
                        <h5>Dune: Messiah</h5>
                        <div style="margin-top: 8px;">
                            <div style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 4px;">Progress: 50%</div>
                            <div style="width: 100%; background: #E5E7EB; height: 6px; border-radius: 3px;">
                                <div style="width: 50%; background: var(--primary); height: 6px; border-radius: 3px;"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="post-actions">
                    <div class="action-btn"><i class="fa-regular fa-heart"></i> 45</div>
                    <div class="action-btn"><i class="fa-regular fa-comment"></i> 2</div>
                    <div class="action-btn"><i class="fa-solid fa-share"></i> Share</div>
                </div>
            </article>

        </main>

        <aside class="widgets">
            
            <div class="widget-card">
                <h3 class="widget-title">Trending Now</h3>
                <ul class="trending-list">
                    <li>
                        <span class="trending-rank">1</span>
                        <div class="trending-info">
                            <strong>Elden Ring DLC</strong>
                            <span>Gaming</span>
                        </div>
                    </li>
                    <li>
                        <span class="trending-rank">2</span>
                        <div class="trending-info">
                            <strong>Shogun (FX)</strong>
                            <span>TV Series</span>
                        </div>
                    </li>
                    <li>
                        <span class="trending-rank">3</span>
                        <div class="trending-info">
                            <strong>Iron Flame</strong>
                            <span>Books</span>
                        </div>
                    </li>
                </ul>
            </div>

            <div class="widget-card">
                <h3 class="widget-title">Friends Online</h3>
                <ul class="trending-list">
                    <li>
                        <img src="https://picsum.photos/seed/friend1/50/50" class="avatar" style="width:30px; height:30px;">
                        <div class="trending-info">
                            <strong>Jessica</strong>
                            <span style="color: var(--accent);">Watching Anime</span>
                        </div>
                    </li>
                    <li>
                        <img src="https://picsum.photos/seed/friend2/50/50" class="avatar" style="width:30px; height:30px;">
                        <div class="trending-info">
                            <strong>Tom</strong>
                            <span style="color: var(--text-muted);">Offline 2m ago</span>
                        </div>
                    </li>
                </ul>
            </div>

            <div class="widget-card" style="background: linear-gradient(135deg, var(--primary), #3730A3);">
                <h3 class="widget-title" style="color: white;">Weekly Challenge</h3>
                <p style="font-size: 0.9rem; color: rgba(255,255,255,0.9); margin-bottom: 12px;">
                    Read 50 pages or watch 3 episodes this week!
                </p>
                <button style="width: 100%; padding: 8px; border-radius: 6px; border: none; background: white; color: var(--primary); font-weight: 600; cursor: pointer;">Join Challenge</button>
            </div>

        </aside>

    </div>

</body>
</html>
```

---

### Part 2: The HTML Template (Dark Mode)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Completionist | Track Everything</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Poppins:wght@600;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    
    <style>
        :root {
            --bg-dark: #111827;
            --bg-card: #1F2937;
            --primary: #6366f1;
            --accent: #10b981;
            --text-main: #F9FAFB;
            --text-muted: #9CA3AF;
            --border: #374151;
            --radius: 12px;
            --nav-width: 260px;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg-dark);
            color: var(--text-main);
            line-height: 1.6;
        }

        a { text-decoration: none; color: inherit; transition: 0.2s; }
        ul { list-style: none; }

        .app-container {
            display: grid;
            grid-template-columns: var(--nav-width) 1fr 320px;
            gap: 24px;
            max-width: 1400px;
            margin: 0 auto;
            min-height: 100vh;
        }

        .sidebar {
            position: sticky;
            top: 0;
            height: 100vh;
            padding: 24px;
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
        }

        .logo {
            font-family: 'Poppins', sans-serif;
            font-size: 1.5rem;
            font-weight: 700;
            color: var(--text-main);
            margin-bottom: 40px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .logo i { color: var(--accent); }

        .nav-menu { display: flex; flex-direction: column; gap: 8px; }
        
        .nav-item {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 12px 16px;
            border-radius: var(--radius);
            color: var(--text-muted);
            font-weight: 500;
        }

        .nav-item:hover, .nav-item.active {
            background-color: rgba(99, 102, 241, 0.1);
            color: var(--primary);
        }

        .nav-item i { width: 20px; text-align: center; }

        .user-mini-profile {
            margin-top: auto;
            display: flex;
            align-items: center;
            gap: 12px;
            padding-top: 20px;
            border-top: 1px solid var(--border);
        }

        .avatar {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            object-fit: cover;
            border: 2px solid var(--primary);
        }

        .user-info h4 { font-size: 0.9rem; color: var(--text-main); }
        .user-info span { font-size: 0.8rem; color: var(--accent); }

        .main-feed {
            padding: 24px 0;
        }

        .feed-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 24px;
        }

        .feed-title { font-family: 'Poppins', sans-serif; font-size: 1.25rem; }
        
        .create-post-btn {
            background-color: var(--primary);
            color: white;
            padding: 10px 20px;
            border-radius: 50px;
            font-weight: 600;
            font-size: 0.9rem;
            border: none;
            cursor: pointer;
        }
        .create-post-btn:hover { opacity: 0.9; }

        .post-card {
            background-color: var(--bg-card);
            border-radius: var(--radius);
            padding: 20px;
            margin-bottom: 20px;
            border: 1px solid var(--border);
        }

        .post-header {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 16px;
        }

        .post-meta h4 { font-size: 0.95rem; }
        .post-meta span { font-size: 0.8rem; color: var(--text-muted); }

        .post-content { margin-bottom: 16px; color: var(--text-main); }

        .completion-item {
            display: flex;
            gap: 16px;
            background: rgba(0,0,0,0.2);
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 16px;
            border-left: 4px solid var(--accent);
        }

        .item-img {
            width: 60px;
            height: 60px;
            border-radius: 6px;
            object-fit: cover;
        }

        .item-details h5 { font-size: 1rem; margin-bottom: 4px; }
        .tag {
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            padding: 2px 6px;
            border-radius: 4px;
            background: #333;
            color: var(--text-muted);
        }
        .tag.anime { color: #F472B6; background: rgba(244, 114, 182, 0.1); }
        .tag.book { color: #60A5FA; background: rgba(96, 165, 250, 0.1); }

        .completion-badge {
            display: flex;
            align-items: center;
            gap: 6px;
            color: var(--accent);
            font-weight: 600;
            font-size: 0.9rem;
            margin-top: 4px;
        }

        .post-actions {
            display: flex;
            gap: 20px;
            border-top: 1px solid var(--border);
            padding-top: 12px;
            color: var(--text-muted);
            font-size: 0.9rem;
        }

        .action-btn { display: flex; align-items: center; gap: 6px; cursor: pointer; }
        .action-btn:hover { color: var(--text-main); }

        .widgets {
            padding: 24px 0;
        }

        .widget-card {
            background-color: var(--bg-card);
            border-radius: var(--radius);
            padding: 20px;
            margin-bottom: 24px;
            border: 1px solid var(--border);
        }

        .widget-title {
            font-family: 'Poppins', sans-serif;
            font-size: 1rem;
            margin-bottom: 16px;
            color: var(--text-main);
        }

        .trending-list li {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 12px;
        }

        .trending-rank {
            font-weight: 700;
            color: var(--text-muted);
            width: 20px;
        }

        .trending-info { font-size: 0.9rem; }
        .trending-info span { display: block; font-size: 0.75rem; color: var(--text-muted); }

        @media (max-width: 1024px) {
            .app-container {
                grid-template-columns: 80px 1fr;
            }
            .sidebar { align-items: center; }
            .nav-item span, .user-info, .logo span { display: none; }
            .logo { margin-bottom: 20px; }
            .user-mini-profile { border-top: none; padding-top: 0; }
            .widgets { display: none; }
        }

        @media (max-width: 600px) {
            .app-container {
                display: block;
            }
            .sidebar {
                position: fixed;
                bottom: 0;
                top: auto;
                width: 100%;
                height: 60px;
                flex-direction: row;
                justify-content: space-around;
                padding: 0;
                background: var(--bg-dark);
                z-index: 100;
                border-top: 1px solid var(--border);
                border-right: none;
            }
            .logo, .user-mini-profile { display: none; }
            .nav-menu { flex-direction: row; width: 100%; justify-content: space-around; }
            .nav-item { padding: 10px; font-size: 1.2rem; }
            .main-feed { padding: 16px; margin-bottom: 60px; }
        }
    </style>
</head>
<body>

    <div class="app-container">
        
        <aside class="sidebar">
            <div class="logo">
                <i class="fa-solid fa-circle-check"></i>
                <span>Completionist</span>
            </div>

            <nav class="nav-menu">
                <a href="#" class="nav-item active">
                    <i class="fa-solid fa-house"></i>
                    <span>Home Feed</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-list-check"></i>
                    <span>My Lists</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-compass"></i>
                    <span>Explore</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-trophy"></i>
                    <span>Challenges</span>
                </a>
                <a href="#" class="nav-item">
                    <i class="fa-solid fa-book-open"></i>
                    <span>Library</span>
                </a>
            </nav>

            <div class="user-mini-profile">
                <img src="https://picsum.photos/seed/me/200/200" alt="User Avatar" class="avatar">
                <div class="user-info">
                    <h4>Alex Doe</h4>
                    <span>Lvl 42 • 500 Items</span>
                </div>
            </div>
        </aside>

        <main class="main-feed">
            
            <header class="feed-header">
                <h2 class="feed-title">Social Feed</h2>
                <button class="create-post-btn"><i class="fa-solid fa-plus"></i> Log Progress</button>
            </header>

            <article class="post-card">
                <div class="post-header">
                    <img src="https://picsum.photos/seed/sarah/200/200" alt="Sarah" class="avatar">
                    <div class="post-meta">
                        <h4>Sarah Jenkins</h4>
                        <span>2 hours ago</span>
                    </div>
                </div>
                
                <div class="post-content">
                    Finally finished this masterpiece! The animation in the last episode was insane. 10/10.
                </div>

                <div class="completion-item">
                    <img src="https://picsum.photos/seed/anime1/200/200" alt="Anime Cover" class="item-img">
                    <div class="item-details">
                        <span class="tag anime">Anime</span>
                        <h5>Neon Cyber Samurai</h5>
                        <div class="completion-badge">
                            <i class="fa-solid fa-circle-check"></i> Completed 12/12 Episodes
                        </div>
                    </div>
                </div>

                <div class="post-actions">
                    <div class="action-btn"><i class="fa-regular fa-heart"></i> 124</div>
                    <div class="action-btn"><i class="fa-regular fa-comment"></i> 18</div>
                    <div class="action-btn"><i class="fa-solid fa-share"></i> Share</div>
                </div>
            </article>

            <article class="post-card">
                <div class="post-header">
                    <img src="https://picsum.photos/seed/mike/200/200" alt="Mike" class="avatar">
                    <div class="post-meta">
                        <h4>Mike Ross</h4>
                        <span>5 hours ago</span>
                    </div>
                </div>
                
                <div class="post-content">
                    Halfway through "Dune". It's getting intense.
                </div>

                <div class="completion-item" style="border-left: 4px solid var(--primary);">
                    <img src="https://picsum.photos/seed/book1/200/200" alt="Book Cover" class="item-img">
                    <div class="item-details">
                        <span class="tag book">Book</span>
                        <h5>Dune: Messiah</h5>
                        <div style="margin-top: 8px;">
                            <div style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 4px;">Progress: 50%</div>
                            <div style="width: 100%; background: #333; height: 6px; border-radius: 3px;">
                                <div style="width: 50%; background: var(--primary); height: 6px; border-radius: 3px;"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="post-actions">
                    <div class="action-btn"><i class="fa-regular fa-heart"></i> 45</div>
                    <div class="action-btn"><i class="fa-regular fa-comment"></i> 2</div>
                    <div class="action-btn"><i class="fa-solid fa-share"></i> Share</div>
                </div>
            </article>

        </main>

        <aside class="widgets">
            
            <div class="widget-card">
                <h3 class="widget-title">Trending Now</h3>
                <ul class="trending-list">
                    <li>
                        <span class="trending-rank">1</span>
                        <div class="trending-info">
                            <strong>Elden Ring DLC</strong>
                            <span>Gaming</span>
                        </div>
                    </li>
                    <li>
                        <span class="trending-rank">2</span>
                        <div class="trending-info">
                            <strong>Shogun (FX)</strong>
                            <span>TV Series</span>
                        </div>
                    </li>
                    <li>
                        <span class="trending-rank">3</span>
                        <div class="trending-info">
                            <strong>Iron Flame</strong>
                            <span>Books</span>
                        </div>
                    </li>
                </ul>
            </div>

            <div class="widget-card">
                <h3 class="widget-title">Friends Online</h3>
                <ul class="trending-list">
                    <li>
                        <img src="https://picsum.photos/seed/friend1/50/50" class="avatar" style="width:30px; height:30px;">
                        <div class="trending-info">
                            <strong>Jessica</strong>
                            <span style="color: var(--accent);">Watching Anime</span>
                        </div>
                    </li>
                    <li>
                        <img src="https://picsum.photos/seed/friend2/50/50" class="avatar" style="width:30px; height:30px;">
                        <div class="trending-info">
                            <strong>Tom</strong>
                            <span style="color: var(--text-muted);">Offline 2m ago</span>
                        </div>
                    </li>
                </ul>
            </div>

            <div class="widget-card" style="background: linear-gradient(135deg, var(--primary), #4338ca);">
                <h3 class="widget-title" style="color: white;">Weekly Challenge</h3>
                <p style="font-size: 0.9rem; color: rgba(255,255,255,0.9); margin-bottom: 12px;">
                    Read 50 pages or watch 3 episodes this week!
                </p>
                <button style="width: 100%; padding: 8px; border-radius: 6px; border: none; background: white; color: var(--primary); font-weight: 600; cursor: pointer;">Join Challenge</button>
            </div>

        </aside>

    </div>

</body>
</html>
```

### Key Features of this Template:
1.  **Social Feed Logic:** The `completion-item` class is designed to handle both "100% Done" items (with a checkmark) and "In-Progress" items (with a progress bar), allowing you to mix different types of content in the feed.
2.  **Visual Hierarchy:**
    *   **Mint Green (`var(--accent)`)** is used sparingly to draw the eye to "Success" and "Completed" states.
    *   **Indigo (`var(--primary)`)** is used for branding and in-progress elements.
3.  **Responsive Design:**
    *   **Desktop:** Shows the full 3-column dashboard.
    *   **Tablet:** Hides the right widget column to give more space to the feed.
    *   **Mobile:** Transforms the left sidebar into a bottom navigation bar (like Instagram or TikTok) for easy thumb access.
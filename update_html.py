import re

html_template = '''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Accountabel AI - Reclaim Your Time</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css2?family=Oswald:wght@700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-dark: #222222;
            --bg-light: #f8f9fa;
            --card-bg: #ffffff;
            --card-border: #e0e0e0;
            --text-light: #ffffff;
            --text-dark: #333333;
            --text-gray: #777777;
            --accent-green: #73B82C;
            --accent-green-hover: #5a9627;
            --danger-color: #ff4d4d;
            --transition-speed: 0.3s;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg-dark);
            color: var(--text-light);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            overflow-x: hidden;
            line-height: 1.6;
        }

        header {
            width: 100%;
            max-width: 1200px;
            margin: 0 auto;
            padding: 1.5rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .logo {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: 1.5rem;
            font-weight: 800;
            color: var(--text-light);
            text-decoration: none;
        }

        .nav-links {
            display: flex;
            gap: 2.5rem;
            list-style: none;
            align-items: center;
        }

        .nav-links a {
            color: var(--text-light);
            text-decoration: none;
            font-weight: 600;
            font-size: 0.95rem;
            transition: opacity var(--transition-speed);
        }

        .nav-links a:hover {
            opacity: 0.8;
        }
        
        .nav-btn {
            background-color: var(--accent-green);
            color: white !important;
            padding: 0.5rem 1.25rem;
            border-radius: 4px;
        }
        .nav-btn:hover {
            background-color: var(--accent-green-hover);
        }

        .hero {
            width: 100%;
            max-width: 1200px;
            margin: 4rem auto;
            padding: 0 2rem;
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 4rem;
            align-items: center;
        }

        .hero-text h1 {
            font-size: 3.5rem;
            font-weight: 800;
            line-height: 1.1;
            margin-bottom: 1.5rem;
        }

        .hero-text p {
            font-size: 1.25rem;
            color: #dddddd;
            margin-bottom: 2.5rem;
            max-width: 90%;
        }

        .cta-group {
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            font-family: inherit;
            font-weight: 700;
            font-size: 1.1rem;
            border-radius: 4px;
            padding: 1rem 2rem;
            cursor: pointer;
            transition: all var(--transition-speed);
            text-decoration: none;
            border: none;
        }

        .btn-primary {
            background-color: var(--accent-green);
            color: white;
        }

        .btn-primary:hover {
            background-color: var(--accent-green-hover);
        }

        .handwritten {
            font-family: 'Comic Sans MS', cursive, sans-serif;
            font-size: 0.9rem;
            transform: rotate(-5deg);
            color: #ddd;
        }

        .hero-image {
            position: relative;
            display: flex;
            justify-content: flex-end;
        }

        .hero-image img {
            max-width: 100%;
            height: auto;
            border-radius: 12px;
            box-shadow: 0 20px 50px rgba(0,0,0,0.5);
        }

        .featured-in {
            text-align: center;
            margin-top: 2rem;
            padding-bottom: 4rem;
        }

        .featured-in span {
            color: #888;
            font-size: 0.9rem;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 1rem;
            display: block;
        }

        .logos {
            display: flex;
            justify-content: center;
            gap: 3rem;
            flex-wrap: wrap;
            opacity: 0.6;
            font-weight: 700;
            font-size: 1.2rem;
            color: white;
            font-family: serif;
        }

        .section-light {
            background-color: var(--bg-light);
            color: var(--text-dark);
            padding: 6rem 2rem;
            text-align: center;
        }

        .section-light h2 {
            font-size: 3rem;
            font-weight: 800;
            text-transform: uppercase;
            margin-bottom: 1rem;
            font-family: 'Oswald', sans-serif;
            letter-spacing: -1px;
        }

        .section-light > p {
            font-size: 1.1rem;
            color: var(--text-gray);
            max-width: 800px;
            margin: 0 auto 4rem;
        }

        .cards {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 2rem;
            max-width: 1200px;
            margin: 0 auto;
        }

        .card {
            background-color: var(--card-bg);
            border-radius: 8px;
            padding: 2rem;
            text-align: left;
            box-shadow: 0 5px 15px rgba(0,0,0,0.05);
            border: 1px solid var(--card-border);
        }

        .card-date {
            font-size: 0.8rem;
            color: var(--text-gray);
            margin-bottom: 0.5rem;
            display: block;
        }

        .card h3 {
            font-size: 1.3rem;
            margin-bottom: 1rem;
            font-weight: 700;
        }

        .card p {
            font-size: 0.95rem;
            color: var(--text-gray);
            margin-bottom: 1.5rem;
        }

        .card-footer {
            display: flex;
            justify-content: space-between;
            font-size: 0.85rem;
            font-weight: 600;
            color: var(--text-dark);
        }

        .features-section {
            background-color: #f0f2f5;
            color: var(--text-dark);
            padding: 6rem 2rem;
            text-align: center;
        }

        .features-section h2 {
            font-size: 2.5rem;
            margin-bottom: 3rem;
            font-weight: 700;
        }

        .checkmarks {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1rem 2rem;
            max-width: 800px;
            margin: 0 auto 3rem;
            text-align: left;
        }

        .check-item {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            font-size: 1.05rem;
            color: #555;
        }

        .check-item svg {
            color: var(--accent-green);
            width: 20px;
            height: 20px;
        }

        .onboarding-section {
            background-color: #ffffff;
            color: var(--text-dark);
            padding: 6rem 2rem;
            text-align: center;
        }
        
        .onboarding-section h2 {
            font-size: 3rem;
            font-weight: 800;
            text-transform: uppercase;
            margin-bottom: 1rem;
            letter-spacing: -1px;
            font-family: 'Oswald', sans-serif;
        }
        
        .onboarding-section > p {
            font-size: 1.1rem;
            color: var(--text-gray);
            margin-bottom: 3rem;
        }

        .onboarding-card {
            background-color: var(--card-bg);
            border: 1px solid var(--card-border);
            border-radius: 12px;
            padding: 2.5rem;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
            max-width: 800px;
            margin: 0 auto;
            text-align: left;
        }

        .steps-indicator {
            display: flex;
            justify-content: space-between;
            margin-bottom: 2.5rem;
            position: relative;
        }

        .steps-indicator::before {
            content: '';
            position: absolute;
            top: 15px;
            left: 0;
            right: 0;
            height: 2px;
            background-color: var(--card-border);
            z-index: 1;
        }

        .step-dot {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            background-color: var(--card-bg);
            border: 2px solid var(--card-border);
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 0.85rem;
            font-weight: 700;
            color: var(--text-gray);
            position: relative;
            z-index: 2;
            transition: all var(--transition-speed);
        }

        .step-dot.active {
            border-color: var(--accent-green);
            color: var(--accent-green);
        }

        .step-dot.completed {
            background-color: var(--accent-green);
            border-color: var(--accent-green);
            color: white;
        }

        .form-step {
            display: none;
            animation: fadeIn 0.4s ease-in-out forwards;
        }

        .form-step.active {
            display: block;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .step-title {
            font-size: 1.75rem;
            font-weight: 700;
            margin-bottom: 0.5rem;
            color: var(--text-dark);
        }

        .step-subtitle {
            font-size: 0.95rem;
            color: var(--text-gray);
            margin-bottom: 2rem;
        }

        .addiction-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 1rem;
            margin-bottom: 2rem;
        }

        .addiction-card {
            background-color: var(--bg-light);
            border: 1px solid var(--card-border);
            border-radius: 8px;
            padding: 1.25rem;
            display: flex;
            align-items: center;
            gap: 1rem;
            cursor: pointer;
            transition: all var(--transition-speed);
            user-select: none;
        }

        .addiction-card:hover {
            border-color: var(--accent-green);
        }

        .addiction-card.selected {
            background-color: rgba(115, 184, 44, 0.1);
            border-color: var(--accent-green);
        }

        .addiction-card .icon {
            font-size: 1.75rem;
        }

        .addiction-card .label {
            font-weight: 600;
            font-size: 1rem;
            color: var(--text-dark);
        }

        .other-input-container {
            grid-column: span 2;
            display: none;
        }

        .other-input-container.active {
            display: block;
        }

        .form-input {
            width: 100%;
            background-color: white;
            border: 1px solid #ccc;
            border-radius: 6px;
            padding: 0.85rem 1.25rem;
            color: var(--text-dark);
            font-family: inherit;
            font-size: 1rem;
            outline: none;
        }

        .form-input:focus {
            border-color: var(--accent-green);
        }

        .input-group {
            margin-bottom: 1.5rem;
        }

        .input-group label {
            display: block;
            font-size: 0.9rem;
            font-weight: 600;
            color: var(--text-gray);
            margin-bottom: 0.5rem;
            text-transform: uppercase;
        }

        .platform-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 1.25rem;
            margin-bottom: 2rem;
        }

        .platform-card {
            background-color: var(--bg-light);
            border: 1px solid var(--card-border);
            border-radius: 12px;
            padding: 2rem 1.5rem;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 1rem;
            cursor: pointer;
            transition: all var(--transition-speed);
            text-align: center;
        }

        .platform-card:hover {
            border-color: var(--accent-green);
        }

        .platform-card.selected {
            background-color: rgba(115, 184, 44, 0.1);
            border-color: var(--accent-green);
        }
        
        .platform-card .icon-img {
            width: 48px;
            height: 48px;
        }

        .platform-card .title {
            font-weight: 700;
            font-size: 1.15rem;
            color: var(--text-dark);
        }

        .btn-container {
            display: flex;
            justify-content: space-between;
            gap: 1rem;
            margin-top: 2.5rem;
        }

        .btn-secondary {
            background-color: white;
            border: 1px solid #ccc;
            color: var(--text-dark);
        }

        .btn-secondary:hover {
            border-color: var(--text-gray);
            background-color: #f5f5f5;
        }

        .btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }

        .code-display-container {
            background-color: var(--bg-light);
            border: 2px dashed var(--accent-green);
            border-radius: 12px;
            padding: 2rem;
            text-align: center;
            margin-bottom: 2rem;
        }

        .code-title {
            font-size: 0.85rem;
            color: var(--text-gray);
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 0.5rem;
        }

        .code-value {
            font-size: 3rem;
            font-weight: 800;
            color: var(--accent-green);
            letter-spacing: 4px;
            font-family: monospace;
            margin-bottom: 0.5rem;
        }

        .code-expiry {
            font-size: 0.8rem;
            color: var(--text-gray);
        }

        .instructions-list {
            list-style: none;
            display: flex;
            flex-direction: column;
            gap: 1.25rem;
            text-align: left;
            margin-bottom: 2rem;
        }

        .instructions-list li {
            display: flex;
            align-items: flex-start;
            gap: 1rem;
        }

        .num-badge {
            width: 24px;
            height: 24px;
            background-color: var(--card-border);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 0.8rem;
            font-weight: 700;
            color: var(--text-dark);
            flex-shrink: 0;
            margin-top: 2px;
        }

        .instruction-content {
            font-size: 0.95rem;
            color: var(--text-gray);
        }

        .instruction-content strong {
            color: var(--text-dark);
        }

        .copy-box {
            background-color: white;
            border: 1px solid var(--card-border);
            border-radius: 6px;
            padding: 0.75rem 1rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 0.5rem;
            font-family: monospace;
            font-size: 0.9rem;
        }

        .copy-btn {
            background: none;
            border: none;
            color: var(--accent-green);
            cursor: pointer;
            font-weight: 700;
            font-family: 'Inter', sans-serif;
            font-size: 0.85rem;
        }

        .copy-btn:hover {
            text-decoration: underline;
        }

        footer {
            background-color: #111;
            color: #888;
            padding: 2rem;
            text-align: center;
            font-size: 0.85rem;
        }

        @media (max-width: 968px) {
            .hero {
                grid-template-columns: 1fr;
                text-align: center;
            }
            .hero-image {
                justify-content: center;
            }
            .cta-group {
                justify-content: center;
            }
            .cards {
                grid-template-columns: 1fr;
            }
            .checkmarks {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>

    <header>
        <a href="#" class="logo">
            Accountabel AI
        </a>
        <ul class="nav-links">
            <li><a href="#onboarding">Products</a></li>
            <li><a href="#onboarding">Features</a></li>
            <li><a href="#onboarding">Support</a></li>
            <li><a href="#onboarding" class="nav-btn">Get Started</a></li>
        </ul>
    </header>

    <main>
        <section class="hero">
            <div class="hero-text">
                <h1>Reclaim your attention.<br>Reclaim your time.</h1>
                <p>Block distractions and let our anti-addiction AI keep you accountable to focus on what actually matters.</p>
                <div class="cta-group">
                    <button class="btn btn-primary" onclick="document.getElementById('onboarding').scrollIntoView({behavior: 'smooth'})">Get for your desktop</button>
                    <span class="handwritten">← free version available!</span>
                </div>
            </div>
            <div class="hero-image">
                <img src="/static/hero_new.png" alt="Accountabel AI Dashboard" onerror="this.src='https://images.unsplash.com/photo-1551288049-bebda4e38f71?auto=format&fit=crop&w=600&q=80'">
            </div>
        </section>

        <div class="featured-in">
            <span>Featured In</span>
            <div class="logos">
                <div>FINANCIAL TIMES</div>
                <div>WIRED</div>
                <div>TIME</div>
                <div>HUFFPOST</div>
                <div>FAST COMPANY</div>
                <div>TechCrunch</div>
                <div>BUSINESS INSIDER</div>
            </div>
        </div>

        <section class="section-light big-tech">
            <h2>BIG TECH ADDICTION</h2>
            <p>It's not just you. Big tech designed their platforms to be addictive. They left us to deal with the consequences. Take back control with Accountabel AI.</p>
            
            <div class="cards">
                <div class="card">
                    <span class="card-date">25 March 2026</span>
                    <h3>Jury in Los Angeles finds Meta and YouTube liable in landmark social media addiction trial</h3>
                    <p>Meta and YouTube must pay millions in damages to a 20-year-old woman after a California jury found the social media giant and video streamer were designed to "hook" young users without concern for their well-being...</p>
                    <div class="card-footer">
                        <span>CBC</span>
                        <span>Source →</span>
                    </div>
                </div>
                <div class="card">
                    <span class="card-date">14 September 2021</span>
                    <h3>A psychiatrist's perspective on social media algorithms and mental health</h3>
                    <p>Variably rewarding users with stimuli keeps them engaged with content. When a user's photo receives a "like," the same dopamine pathways involved in motivation, reward, and addiction are activated...</p>
                    <div class="card-footer">
                        <span>Stanford</span>
                        <span>Source →</span>
                    </div>
                </div>
                <div class="card">
                    <span class="card-date">20 February 2019</span>
                    <h3>Cognitive deficits in problematic internet use: meta-analysis of 40 studies</h3>
                    <p>In meta-analysis, [problematic internet use] was associated with significant cognitive deficits in attentional inhibition, motor inhibition (and pre-potent motor inhibition), decision-making and working memory...</p>
                    <div class="card-footer">
                        <span>Cambridge</span>
                        <span>Source →</span>
                    </div>
                </div>
            </div>
        </section>

        <section class="features-section">
            <h2>Flexible features for any goal.</h2>
            <div class="checkmarks">
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Block or allow websites
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Set weekly schedules
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Block apps and games
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Set allowances and other breaks
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Block adult websites
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Track website and app usage
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Personalized AI Check-ins
                </div>
                <div class="check-item">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    Craving Analysis & Support
                </div>
            </div>
            <button class="btn btn-primary" onclick="document.getElementById('onboarding').scrollIntoView({behavior: 'smooth'})">See the full set of features</button>
        </section>

        <section id="onboarding" class="onboarding-section">
            <h2>JOIN OTHERS WINNING TIME BACK</h2>
            <p>Downloaded millions of times, saving lifetimes of wasted time.</p>
            
            <div class="onboarding-card">
                <!-- Progress Steps Indicator -->
                <div class="steps-indicator">
                    <div class="step-dot active" id="dot-1">1</div>
                    <div class="step-dot" id="dot-2">2</div>
                    <div class="step-dot" id="dot-3">3</div>
                    <div class="step-dot" id="dot-4">4</div>
                </div>

                <form id="onboarding-form" onsubmit="event.preventDefault();">
                    <!-- Step 1: Addiction/Habit selection -->
                    <div class="form-step active" id="step-1">
                        <h2 class="step-title">Choose Your Habits</h2>
                        <p class="step-subtitle">Select the habits or addictions you want accountability for. You can choose multiple.</p>
                        
                        <div class="addiction-grid">
                            <div class="addiction-card" data-val="Social Media">
                                <span class="icon">📱</span>
                                <span class="label">Social Media</span>
                            </div>
                            <div class="addiction-card" data-val="Alcohol">
                                <span class="icon">🍷</span>
                                <span class="label">Alcohol</span>
                            </div>
                            <div class="addiction-card" data-val="Adult Content">
                                <span class="icon">🔞</span>
                                <span class="label">Adult Content</span>
                            </div>
                            <div class="addiction-card" data-val="Cigarettes">
                                <span class="icon">🚬</span>
                                <span class="label">Cigarettes</span>
                            </div>
                            <div class="addiction-card" data-val="Drugs">
                                <span class="icon">💊</span>
                                <span class="label">Drugs</span>
                            </div>
                            <div class="addiction-card" data-val="Gaming">
                                <span class="icon">🎮</span>
                                <span class="label">Gaming</span>
                            </div>
                            <div class="addiction-card" data-val="Shopping">
                                <span class="icon">🛍️</span>
                                <span class="label">Shopping</span>
                            </div>
                            <div class="addiction-card" id="card-other" data-val="Other">
                                <span class="icon">✏️</span>
                                <span class="label">Custom Habit</span>
                            </div>

                            <div class="other-input-container" id="custom-habit-container">
                                <input type="text" id="custom-habit-input" class="form-input" placeholder="Enter custom habit (e.g. junk food, caffeine)">
                            </div>
                        </div>
                    </div>

                    <!-- Step 2: Timezone & Check-in time -->
                    <div class="form-step" id="step-2">
                        <h2 class="step-title">Daily Check-In Schedule</h2>
                        <p class="step-subtitle">Decide when the accountability AI bot should message you to check in on your craving level.</p>

                        <div class="input-group">
                            <label for="timezone-select">Select Your Timezone</label>
                            <select id="timezone-select" class="form-input" style="appearance: auto;">
                                <option value="UTC">UTC (GMT+00:00)</option>
                                <option value="EST">EST (GMT-05:00)</option>
                                <option value="CST">CST (GMT-06:00)</option>
                                <option value="PST">PST (GMT-08:00)</option>
                                <option value="GMT">GMT (GMT+00:00)</option>
                                <option value="BST">BST (GMT+01:00)</option>
                                <option value="IST">IST (GMT+05:30)</option>
                                <option value="AEST">AEST (GMT+10:00)</option>
                            </select>
                        </div>

                        <div class="input-group">
                            <label for="checkin-time">Daily Check-In Time (24h format)</label>
                            <input type="time" id="checkin-time" class="form-input" value="21:00" required>
                        </div>
                    </div>

                    <!-- Step 3: Platform selection -->
                    <div class="form-step" id="step-3">
                        <h2 class="step-title">Choose Platform</h2>
                        <p class="step-subtitle">Where do you want to interact with your accountability companion?</p>

                        <div class="platform-grid">
                            <div class="platform-card selected" data-platform="telegram">
                                <svg class="icon-img" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12C22 6.48 17.52 2 12 2ZM16.64 8.8C16.49 10.38 15.84 14.18 15.51 15.96C15.37 16.71 15.09 16.96 14.82 16.99C14.24 17.04 13.8 16.61 13.24 16.24C12.36 15.66 11.86 15.3 11.01 14.74C10.02 14.09 10.66 13.73 11.23 13.14C11.38 12.99 13.95 10.66 14 10.45C14.0063 10.4079 14.0016 10.3644 13.9863 10.3248C13.971 10.2852 13.9458 10.2514 13.9137 10.2275C13.84 10.17 13.73 10.19 13.65 10.21C13.54 10.23 11.82 11.37 8.49 13.62C8 13.96 7.55 14.12 7.15 14.11C6.71 14.1 5.86 13.86 5.23 13.66C4.46 13.41 3.85 13.28 3.9 12.86C3.93 12.64 4.23 12.41 4.8 12.17C8.33 10.63 10.7 9.61 11.89 9.12C15.31 7.7 16.01 7.45 16.48 7.46C16.58 7.46 16.8 7.48 16.94 7.6C17.06 7.7 17.09 7.86 17.11 7.97C17.11 8.04 17.13 8.35 17.11 8.56C16.94 8.78 16.64 8.8 16.64 8.8Z" fill="%2373B82C"/>
                                </svg>
                                <span class="title">Telegram</span>
                            </div>
                            <div class="platform-card" data-platform="discord">
                                <svg class="icon-img" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M19.27 4.73C17.8 4.05 16.2 3.56 14.53 3.3C14.32 3.67 14.09 4.14 13.93 4.53C12.16 4.26 10.4 4.26 8.66 4.53C8.5 4.14 8.26 3.67 8.05 3.3C6.38 3.56 4.78 4.05 3.31 4.73C0.34 9.17 -0.46 13.51 0.17 17.79C2.14 19.24 4.05 20.12 5.92 20.7C6.38 20.07 6.79 19.39 7.14 18.67C6.47 18.42 5.83 18.11 5.23 17.75C5.39 17.63 5.55 17.51 5.7 17.38C9.52 19.14 13.68 19.14 17.46 17.38C17.61 17.51 17.77 17.63 17.93 17.75C17.33 18.11 16.69 18.42 16.02 18.67C16.37 19.39 16.78 20.07 17.24 20.7C19.11 20.12 21.02 19.24 22.99 17.79C23.76 12.4 21.78 8.1 19.27 4.73ZM8.14 15C6.98 15 6.03 13.93 6.03 12.62C6.03 11.31 6.96 10.24 8.14 10.24C9.32 10.24 10.28 11.31 10.25 12.62C10.25 13.93 9.32 15 8.14 15ZM15.06 15C13.9 15 12.95 13.93 12.95 12.62C12.95 11.31 13.88 10.24 15.06 10.24C16.24 10.24 17.2 11.31 17.17 12.62C17.17 13.93 16.24 15 15.06 15Z" fill="%2373B82C"/>
                                </svg>
                                <span class="title">Discord</span>
                            </div>
                        </div>
                    </div>

                    <!-- Step 4: Connecting code & instructions -->
                    <div class="form-step" id="step-4">
                        <h2 class="step-title">Link Your Account</h2>
                        <p class="step-subtitle">To pair your settings, copy the code below and send it to the accountability bot.</p>

                        <div class="code-display-container">
                            <div class="code-title">Your Pairing Code</div>
                            <div class="code-value" id="pairing-code">------</div>
                            <div class="code-expiry">Expires in 15 minutes</div>
                        </div>

                        <ul class="instructions-list">
                            <li>
                                <div class="num-badge">1</div>
                                <div class="instruction-content">
                                    Open the bot on <strong id="platform-text-instruction">Telegram</strong>:
                                    <br>
                                    <a href="#" id="bot-redirect-link" class="btn btn-secondary" style="margin-top:0.5rem; padding: 0.5rem 1rem; display:inline-flex;" target="_blank">Open Bot Chat ↗</a>
                                </div>
                            </li>
                            <li>
                                <div class="num-badge">2</div>
                                <div class="instruction-content">
                                    Copy and send the pairing command directly to the bot chat:
                                    <div class="copy-box">
                                        <span id="command-text">/link ------</span>
                                        <button type="button" class="copy-btn" id="btn-copy-cmd" onclick="copyCommand()">COPY</button>
                                    </div>
                                </div>
                            </li>
                        </ul>
                    </div>

                    <!-- Step navigation buttons -->
                    <div class="btn-container">
                        <button type="button" class="btn btn-secondary" id="btn-prev" onclick="changeStep(-1)" disabled>Back</button>
                        <button type="button" class="btn btn-primary" id="btn-next" onclick="changeStep(1)">Next</button>
                    </div>
                </form>
            </div>
        </section>

    </main>

    <footer>
        <div>&copy; 2026 Accountabel AI. All rights reserved.</div>
    </footer>
'''

with open("static/index.html", "r", encoding="utf-8") as f:
    old_html = f.read()

# Extract script logic
script_match = re.search(r'<script>(.*?)</script>', old_html, re.DOTALL)
if script_match:
    js_code = script_match.group(1)
else:
    print("Script not found!")
    exit(1)

# we can inject the script at the end
final_html = html_template + "\n<script>\n" + js_code + "\n</script>\n</body>\n</html>"

with open("static/index.html", "w", encoding="utf-8") as f:
    f.write(final_html)

print("Updated index.html successfully.")

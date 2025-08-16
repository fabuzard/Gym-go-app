# 🏋️ Gym API – 2025  

A **RESTful API** for a gym and fitness tracking platform where users can register, log in, create workouts, manage exercises, and track performance logs. Built with **Go** using the **Echo framework**, the system applies **REST best practices**, JWT-based authentication, and ownership validation for secure resource access.  

The project integrates a third-party **BMI Calculator API** to provide personalized health insights, supports **Swagger documentation**, and is deployed on **Heroku** with environment-based configuration.  

## 🚀 Features  

- **User Management**  
  - Register, log in, and manage user accounts  
  - JWT-based authentication and secure password handling  
  - Ownership validation for all resources  

- **Workouts & Exercises**  
  - Create, update, delete, and view workouts  
  - Add exercises to workouts and track logs  
  - Cascade delete for workouts and exercises  

- **Performance Tracking**  
  - Log exercise sets, reps, weights, and timestamps  
  - Retrieve all logs for authenticated users  

- **Third-Party Integration**  
  - Integrated **RapidAPI BMI Calculator** to calculate BMI from user’s profile data  

- **API Documentation**  
  - Swagger UI available at `/swagger/index.html`  

- **Deployment**  
  - Live API deployed on **Heroku**  
  - Database hosted on **Supabase (PostgreSQL)**  

## 🛠️ Tech Stack  

- **Backend:** Go, Echo  
- **Database:** PostgreSQL (Supabase)  
- **Auth:** JWT (JSON Web Tokens)  
- **Docs:** Swagger  
- **Hosting:** Heroku  
- **Third-Party:** RapidAPI (BMI Calculator)  

## 📌 API Endpoints  

### 🔑 Authentication  
- **POST** `/api/users/register` → Register new user  
- **POST** `/api/users/login` → Login and get JWT token  

### 👤 User  
- **GET** `/api/users` → Get authenticated user info + BMI (via third-party API)  

### 🏋️ Workouts  
- **GET** `/api/workouts` → List all workouts (owned by user)  
- **GET** `/api/workouts/:id` → Get workout + related exercises  
- **POST** `/api/workouts` → Create new workout  
- **PUT** `/api/workouts/:id` → Update workout (owner-only)  
- **DELETE** `/api/workouts/:id` → Delete workout + exercises (owner-only)  

### 🏃 Exercises  
- **POST** `/api/exercises` → Add exercise to a workout  
- **DELETE** `/api/exercises/:id` → Delete exercise (owner-only)  

### 📊 Logs  
- **POST** `/api/logs` → Create exercise log (weights, reps, sets)  
- **GET** `/api/logs` → Get all logs for authenticated user  

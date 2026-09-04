# Go-Project

یک پروژه‌ی **بک‌اند API** نوشته‌شده با زبان Go، با ساختاری تمیز و قابل توسعه که برای پیاده‌سازی سرویس‌های واقعی و آماده‌ی پروداکشن طراحی شده است.

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-1.2x-00ADD8?logo=go&logoColor=white">
  <img alt="Gin" src="https://img.shields.io/badge/Framework-Gin-00ADD8">
  <img alt="GORM" src="https://img.shields.io/badge/ORM-GORM-informational">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/Database-PostgreSQL-336791?logo=postgresql&logoColor=white">
  <img alt="Redis" src="https://img.shields.io/badge/Cache-Redis-DC382D?logo=redis&logoColor=white">
  <img alt="Elasticsearch" src="https://img.shields.io/badge/Search-Elasticsearch-005571?logo=elasticsearch&logoColor=white">
  <img alt="License" src="https://img.shields.io/badge/License-MIT-green">
</p>

---

## 📋 درباره‌ی پروژه

این ریپازیتوری فقط شامل بخش **API** است و با هدف ایجاد یک اسکلت (Boilerplate) حرفه‌ای و قابل استفاده‌ی مجدد برای پروژه‌های مبتنی بر Go ساخته شده است. معماری آن به‌گونه‌ای طراحی شده که افزودن سرویس، مدل و اندپوینت جدید ساده و بدون درهم‌ریختگی باشد.

## ✨ امکانات و قابلیت‌ها

- 🔐 **احراز هویت با JWT** به همراه میدلویر‌های مخصوص کنترل دسترسی
- 📱 **سیستم OTP** برای ورود/ثبت‌نام امن با کد یک‌بار مصرف
- 👥 **مدیریت کاربران** شامل ثبت‌نام، پروفایل و سطوح دسترسی (نقش‌محور / Role-Based Access Control)
- 🏙️ **مدیریت شهرها** (ایجاد، ویرایش، دریافت لیست و ...)
- 📁 **آپلود و مدیریت فایل**
- 🗄️ **پایگاه داده PostgreSQL** با ORM قدرتمند **GORM**
- ⚡ **کش با Redis** برای افزایش سرعت پاسخ‌گویی
- 🔍 **جست‌وجوی متنی با Elasticsearch**
- 🧱 **معماری لایه‌ای و ماژولار** (Handler / Service / Repository)
- 🌐 **فریم‌ورک Gin** برای مسیریابی و مدیریت درخواست‌ها

## 🛠️ تکنولوژی‌های استفاده‌شده

| بخش | تکنولوژی |
|---|---|
| زبان برنامه‌نویسی | Go |
| وب فریم‌ورک | Gin |
| ORM | GORM |
| پایگاه داده | PostgreSQL |
| کش | Redis |
| موتور جست‌وجو | Elasticsearch |
| احراز هویت | JWT |


## 🔐 احراز هویت

برای دسترسی به اندپوینت‌های محافظت‌شده، توکن JWT دریافتی از مرحله‌ی ورود را در هدر درخواست ارسال کنید:

```
Authorization: Bearer <your_token>
```


<p align="center">Made with ❤️ using Go</p>

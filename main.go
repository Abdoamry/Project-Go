package main 

import(
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/yourname/fiber-jwt-app/config"
	"github.com/yourname/fiber-jwt-app/database"
	"github.com/yourname/fiber-jwt-app/routes"
	"github.com/yourname/fiber-jwt-app/utils/logger"

	"github.com/gofiber/fiber/v2"
)



func main() {
	godotenv.Load() // لتحميل المتغيرات من ملف .env

	logger.InitLogger() // إعداد zap
	database.ConnectDB() // الاتصال بقاعدة البيانات
	database.InitRedis() // إعداد Redis

	app := fiber.New()

	routes.SetupRoutes(app) // تعريف المسارات

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}



// 📌 **ملاحظات الربط:**
// - Fiber = السيرفر الأساسي + إدارة الـ routes.
// - GORM = ORM للربط بقاعدة PostgreSQL.
// - godotenv = لتحميل الإعدادات من ملف .env.
// - zap = لعمل logging احترافي.
// - Redis = لتخزين بيانات مؤقتة (مثل JWT token).
// - JWT = للتوثيق وتسجيل الدخول.

// كل جزء متكامل ويشتغل مع الباقي كالتالي:
// - المستخدم يسجل أو يسجل دخول → يصدر له توكن JWT.
// - التوكن يُخزن في Redis + يُستخدم للتحقق لاحقًا.
// - Fiber و middleware بيمنع الوصول للراوتات المحمية بدون JWT.
// - قاعدة البيانات تحفظ بيانات المستخدمين.
// - zap بيسجل logs التشغيل المهمة.
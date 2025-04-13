### **解决方案**  
需要修改后端登录接口的响应内容，**在返回 JWT Token 的同时，显式返回用户角色（`role`）**，便于前端直接根据角色跳转页面。以下是具体实现步骤：

---

### **1. 后端代码修改（`auth.go`）**  
在 `Login` 函数中，将用户角色添加到响应 JSON 中：  
```go
func Login(c *fiber.Ctx) error {
    // ...（原有代码，查询用户并验证密码）

    // 生成 JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role, // JWT 中已包含 role
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })
    tokenString, err := token.SignedString([]byte(config.Get("JWT_SECRET")))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
    }

    // 返回 Token 和 Role
    return c.JSON(fiber.Map{
        "token": tokenString,
        "role":  user.Role, // 显式返回 role
    })
}
```

---

### **2. 前端处理登录响应（Vue + Axios）**  
在登录成功后，前端从响应中提取 `role` 并存储到本地（如 Vuex 或 localStorage），然后根据角色跳转页面：  
```javascript
// 示例：Vue 组件中的登录方法
import axios from 'axios';

export default {
  methods: {
    async handleLogin() {
      try {
        const response = await axios.post('/api/login', {
          username: this.username,
          password: this.password,
        });
        
        // 存储 Token 和 Role
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('role', response.data.role);
        
        // 根据角色跳转不同页面
        switch (response.data.role) {
          case 'seller':
            this.$router.push('/seller-dashboard');
            break;
          case 'client':
            this.$router.push('/client-dashboard');
            break;
          case 'sales':
            this.$router.push('/sales-dashboard');
            break;
          default:
            this.$router.push('/');
        }
      } catch (error) {
        console.error('登录失败:', error.response.data.error);
      }
    },
  },
};
```

---

### **3. 前端路由守卫（权限控制）**  
在 Vue Router 中配置全局前置守卫，根据本地存储的 `role` 控制访问权限：  
```javascript
// router.js
import Vue from 'vue';
import Router from 'vue-router';
import store from './store';

Vue.use(Router);

const router = new Router({
  routes: [
    { path: '/seller-dashboard', component: SellerDashboard, meta: { requiresRole: 'seller' } },
    { path: '/client-dashboard', component: ClientDashboard, meta: { requiresRole: 'client' } },
    { path: '/sales-dashboard', component: SalesDashboard, meta: { requiresRole: 'sales' } },
  ],
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role');

  // 检查路由是否需要特定角色
  if (to.meta.requiresRole) {
    if (!token || role !== to.meta.requiresRole) {
      next('/login'); // 无权限则跳转登录页
      return;
    }
  }
  next();
});
```

---

### **4. 后端安全性增强（可选）**  
如果前端需要更严格的权限控制，可以在关键接口（如 `/api/seller/orders`）中，后端校验 JWT 中的 `role` 字段，确保用户角色与接口权限匹配。  
（已在你的中间件 `auth.go` 中实现角色校验，无需额外修改）

---

### **最终效果**  
1. **登录流程**：  
   - 用户登录 → 后端返回 `token` 和 `role` → 前端存储并跳转对应页面。  
2. **权限控制**：  
   - 未登录用户访问受限路由 → 自动跳转登录页。  
   - 普通用户尝试访问管理员页面 → 自动拦截。  
3. **前后端协作**：  
   - 前端无需解码 JWT，直接使用显式返回的 `role`，简化逻辑。  
   - 后端通过中间件确保接口安全性。  

通过以上调整，系统将实现**基于角色的动态路由跳转和权限控制**，完全符合你的需求。
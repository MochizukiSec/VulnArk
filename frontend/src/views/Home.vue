<template>
  <div class="home-container">
    <!-- 导航栏 -->
    <header class="home-header" :class="{ 'header-hidden': !isHeaderVisible }">
      <div class="header-container">
        <div class="logo">
          <div class="logo-icon">
            <div class="logo-animation">
              <div class="logo-shield">
                <div class="logo-shield-inner"></div>
                <div class="logo-sparkles">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
              <div class="logo-pulse"></div>
            </div>
          </div>
          <span class="logo-text">VulnArk</span>
        </div>
        
        <nav class="main-nav">
          <ul class="nav-links">
            <li><a href="#features" class="nav-link" :class="{ 'active': activeSection === 'features' }">功能</a></li>
            <li><a href="#workflow" class="nav-link" :class="{ 'active': activeSection === 'workflow' }">工作流程</a></li>
            <li><a href="#faq" class="nav-link" :class="{ 'active': activeSection === 'faq' }">常见问题</a></li>
          </ul>
        </nav>
        
        <div class="header-actions">
          <router-link to="/login" class="header-btn login-btn">登录</router-link>
        </div>
        
        <button class="mobile-menu-btn" @click="toggleMobileMenu">
          <i class="el-icon-menu"></i>
        </button>
      </div>
    </header>

    <!-- 移动端菜单 -->
    <div class="mobile-menu" :class="{ 'mobile-menu-open': mobileMenuOpen }">
      <div class="mobile-menu-content">
        <div class="mobile-menu-header">
          <div class="logo">
            <div class="logo-icon">
              <div class="logo-animation">
                <div class="logo-shield">
                  <div class="logo-shield-inner"></div>
                  <div class="logo-sparkles">
                    <span></span>
                    <span></span>
                    <span></span>
                  </div>
                </div>
                <div class="logo-pulse"></div>
              </div>
            </div>
            <span class="logo-text">VulnArk</span>
          </div>
          <button class="mobile-menu-close" @click="toggleMobileMenu">
            <i class="el-icon-close"></i>
          </button>
        </div>
        <ul class="mobile-nav-links">
          <li><a href="#features" class="mobile-nav-link" :class="{ 'active': activeSection === 'features' }" @click="toggleMobileMenu">功能</a></li>
          <li><a href="#workflow" class="mobile-nav-link" :class="{ 'active': activeSection === 'workflow' }" @click="toggleMobileMenu">工作流程</a></li>
          <li><a href="#faq" class="mobile-nav-link" :class="{ 'active': activeSection === 'faq' }" @click="toggleMobileMenu">常见问题</a></li>
        </ul>
        <div class="mobile-menu-actions">
          <router-link to="/login" class="header-btn login-btn" @click="toggleMobileMenu">登录</router-link>
        </div>
      </div>
    </div>

    <!-- 背景装饰 -->
    <div class="background-decoration">
      <!-- 背景图层 -->
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
      <div class="bg-dots"></div>
    </div>

    <div class="home-content">
      <!-- 英雄区域 -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">VulnArk</h1>
          <p class="hero-subtitle">安全漏洞全生命周期管理系统，集成报告生成与数据分析功能</p>
          <div class="hero-actions">
            <router-link to="/login" class="hero-btn primary-btn">
              <span>开始使用</span>
              <i class="el-icon-arrow-right"></i>
            </router-link>
            <a href="#features" class="hero-btn secondary-btn">
              <span>了解功能</span>
              <i class="el-icon-d-arrow-right"></i>
            </a>
          </div>
        </div>
        
        <!-- 添加动画图像 -->
        <div class="hero-image">
          <div class="security-animation">
            <!-- 安全动画或图表 -->
            <div class="security-shield">
              <i class="el-icon-s-check shield-icon"></i>
              <div class="shield-pulse"></div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 特性部分 -->
      <div id="features" class="features-section">
        <h2 class="section-title">主要功能</h2>
        <div class="section-subtitle">全面的安全漏洞管理平台，提供完整的安全生命周期管理</div>
        
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">
              <i class="el-icon-warning-outline"></i>
            </div>
            <h3>漏洞管理</h3>
            <p>全面管理系统安全漏洞，基于CVSS评分评估风险</p>
            <ul class="feature-list">
              <li><span class="feature-highlight">多维度分类</span> 支持多种漏洞属性</li>
              <li><span class="feature-highlight">生命周期管理</span> 从发现到修复的完整跟踪</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="el-icon-data-analysis"></i>
            </div>
            <h3>仪表盘分析</h3>
            <p>直观的安全数据可视化，帮助快速识别安全趋势</p>
            <ul class="feature-list">
              <li><span class="feature-highlight">实时统计</span> 关键安全指标一目了然</li>
              <li><span class="feature-highlight">分布分析</span> 漏洞分布和修复进度跟踪</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="el-icon-document"></i>
            </div>
            <h3>安全报告</h3>
            <p>强大的报告生成功能，支持多种报告类型和格式</p>
            <ul class="feature-list">
              <li><span class="feature-highlight">多种报告类型</span> 摘要、详细、合规和趋势分析</li>
              <li><span class="feature-highlight">多格式导出</span> 支持PDF、Excel、Word和HTML</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="el-icon-user"></i>
            </div>
            <h3>用户与权限</h3>
            <p>完善的用户管理和权限控制系统</p>
            <ul class="feature-list">
              <li><span class="feature-highlight">角色权限</span> 管理员、普通用户和只读用户</li>
              <li><span class="feature-highlight">用户状态管理</span> 活跃、非活跃和禁用状态</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="el-icon-s-tools"></i>
            </div>
            <h3>系统集成</h3>
            <p>灵活的API设计，便于与其他安全工具集成</p>
            <ul class="feature-list">
              <li><span class="feature-highlight">RESTful API</span> 标准化的接口设计</li>
              <li><span class="feature-highlight">批量导入</span> 支持批量漏洞数据导入</li>
            </ul>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="el-icon-lock"></i>
            </div>
            <h3>安全与合规</h3>
            <p>内置的安全机制，确保系统合规和安全</p>
            <ul class="feature-list">
              <li><span class="feature-highlight">JWT认证</span> 安全的身份验证机制</li>
              <li><span class="feature-highlight">审计日志</span> 关键操作的审计跟踪</li>
            </ul>
          </div>
        </div>
      </div>
      
      <!-- 添加工作流程部分 -->
      <div id="workflow" class="workflow-section">
        <h2 class="section-title">工作流程</h2>
        <div class="section-subtitle">系统化的漏洞管理流程，实现高效的安全问题处理</div>
        
        <div class="workflow-steps">
          <div class="workflow-step">
            <div class="step-number">1</div>
            <div class="step-content">
              <h3>漏洞记录</h3>
              <p>通过界面录入或API批量导入安全漏洞信息</p>
              <div class="step-icon">
                <i class="el-icon-edit-outline"></i>
              </div>
            </div>
          </div>
          <div class="workflow-step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h3>风险评估</h3>
              <p>使用CVSS评分系统评估漏洞风险等级并确定优先级</p>
              <div class="step-icon">
                <i class="el-icon-rank"></i>
              </div>
            </div>
          </div>
          <div class="workflow-step">
            <div class="step-number">3</div>
            <div class="step-content">
              <h3>状态跟踪</h3>
              <p>跟踪漏洞从开放到解决的全过程，记录所有处理步骤</p>
              <div class="step-icon">
                <i class="el-icon-s-promotion"></i>
              </div>
            </div>
          </div>
          <div class="workflow-step">
            <div class="step-number">4</div>
            <div class="step-content">
              <h3>数据分析</h3>
              <p>通过仪表盘查看安全态势，识别风险趋势</p>
              <div class="step-icon">
                <i class="el-icon-pie-chart"></i>
              </div>
            </div>
          </div>
          <div class="workflow-step">
            <div class="step-number">5</div>
            <div class="step-content">
              <h3>报告生成</h3>
              <p>生成多种类型的安全报告，支持多种文档格式</p>
              <div class="step-icon">
                <i class="el-icon-document-checked"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 添加常见问题部分 -->
      <div id="faq" class="faq-section">
        <h2 class="section-title">常见问题</h2>
        <div class="section-subtitle">关于我们漏洞管理平台的常见疑问解答</div>
        
        <div class="faq-container">
          <div class="faq-item">
            <div class="faq-question">
              <h3>系统支持哪些类型的用户角色？</h3>
              <i class="el-icon-arrow-down"></i>
            </div>
            <div class="faq-answer">
              <p>我们的系统支持三种用户角色：管理员、普通用户和只读用户。管理员拥有所有权限，可以管理用户和系统配置；普通用户可以创建和管理漏洞以及生成报告；只读用户仅能查看信息但不能修改数据。</p>
            </div>
          </div>
          <div class="faq-item">
            <div class="faq-question">
              <h3>平台支持哪些报告格式？</h3>
              <i class="el-icon-arrow-down"></i>
            </div>
            <div class="faq-answer">
              <p>我们的平台支持多种报告格式，包括PDF、Excel、Word和HTML。可以根据需要生成摘要报告、详细报告、合规报告和趋势分析报告，满足不同场景下的安全报告需求。</p>
            </div>
          </div>
          <div class="faq-item">
            <div class="faq-question">
              <h3>系统如何帮助跟踪漏洞生命周期？</h3>
              <i class="el-icon-arrow-down"></i>
            </div>
            <div class="faq-answer">
              <p>系统提供完整的漏洞状态管理功能，从开放、处理中、已解决到已关闭，记录每个状态的变更历史和相关信息。同时支持根据漏洞严重程度（严重、高危、中危、低危和信息）进行优先级管理和资源分配。</p>
            </div>
          </div>
          <div class="faq-item">
            <div class="faq-question">
              <h3>平台的数据分析功能有哪些？</h3>
              <i class="el-icon-arrow-down"></i>
            </div>
            <div class="faq-answer">
              <p>我们的平台提供实时仪表盘，展示漏洞分布、解决率、风险等级分布等关键安全指标。支持根据时间、部门、系统等多维度进行数据分析，帮助安全团队了解整体安全态势并制定相应策略。</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 行动召唤部分 -->
      <div class="cta-section">
        <div class="cta-content">
          <h2>准备好提升您的安全管理了吗？</h2>
          <p>立即开始使用我们的漏洞管理平台</p>
          <div class="cta-actions">
            <router-link to="/login" class="hero-btn primary-btn">
              <span>开始使用</span>
              <i class="el-icon-arrow-right"></i>
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- 添加页脚 -->
      <footer class="home-footer">
        <div class="footer-content">
          <div class="footer-grid">
            <div class="footer-col">
              <div class="footer-logo">
                <div class="logo">
                  <div class="logo-icon">
                    <div class="logo-animation">
                      <div class="logo-shield">
                        <div class="logo-shield-inner"></div>
                        <div class="logo-sparkles">
                          <span></span>
                          <span></span>
                          <span></span>
                        </div>
                      </div>
                      <div class="logo-pulse"></div>
                    </div>
                  </div>
                  <span class="logo-text">VulnArk</span>
                </div>
              </div>
              <p class="footer-desc">打造业界领先的安全漏洞管理解决方案</p>
            </div>
            <div class="footer-col">
              <h3>产品</h3>
              <ul class="footer-links">
                <li><a href="#features">功能特性</a></li>
                <li><a href="#workflow">工作流程</a></li>
                <li><a href="#faq">常见问题</a></li>
              </ul>
            </div>
            <div class="footer-col">
              <h3>资源</h3>
              <ul class="footer-links">
                <li><a href="#">安全公告</a></li>
                <li><a href="#">最佳实践</a></li>
                <li><a href="#">漏洞百科</a></li>
                <li><a href="#">安全论坛</a></li>
              </ul>
            </div>
          </div>
          <div class="footer-bottom">
            <p>&copy; 2023 VulnArk. 保留所有权利.</p>
            <div class="footer-social">
              <a href="#" class="social-link"><i class="el-icon-s-platform"></i></a>
              <a href="#" class="social-link"><i class="el-icon-s-custom"></i></a>
              <a href="#" class="social-link"><i class="el-icon-s-promotion"></i></a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  </div>
</template>

<script>
import { onMounted, ref, onUnmounted } from 'vue'

export default {
  name: 'Home',
  
  setup() {
    const activeFaq = ref(null);
    const activeSection = ref('');
    const mobileMenuOpen = ref(false);
    const isHeaderVisible = ref(true);
    const lastScrollTop = ref(0);
    
    const toggleFaq = (index) => {
      activeFaq.value = activeFaq.value === index ? null : index;
    };
    
    const toggleMobileMenu = () => {
      mobileMenuOpen.value = !mobileMenuOpen.value;
      if (mobileMenuOpen.value) {
        document.body.style.overflow = 'hidden'; // 防止菜单打开时背景滚动
      } else {
        document.body.style.overflow = 'auto';
      }
    };
    
    const handleScroll = () => {
      const currentScrollTop = window.scrollY;
      // 判断滚动方向
      if (currentScrollTop > lastScrollTop.value && currentScrollTop > 80) {
        // 向下滚动，隐藏导航栏
        isHeaderVisible.value = false;
      } else {
        // 向上滚动，显示导航栏
        isHeaderVisible.value = true;
      }
      lastScrollTop.value = currentScrollTop;
      
      // 检查当前活动区域
      checkScrollPosition();
    };
    
    const checkScrollPosition = () => {
      // 获取各部分位置
      const featuresEl = document.getElementById('features');
      const workflowEl = document.getElementById('workflow');
      const faqEl = document.getElementById('faq');
      
      if (!featuresEl || !workflowEl || !faqEl) return;
      
      // 获取当前滚动位置，添加一个偏移量使高亮更加准确
      const scrollPosition = window.scrollY;
      const offset = 300; // 调整这个值可以更早地高亮下一个区域
      
      // 简单直接地判断当前区域
      if (scrollPosition + offset >= faqEl.offsetTop) {
        activeSection.value = 'faq';
      } else if (scrollPosition + offset >= workflowEl.offsetTop) {
        activeSection.value = 'workflow';
      } else if (scrollPosition + offset >= featuresEl.offsetTop) {
        activeSection.value = 'features';
      } else {
        activeSection.value = '';
      }
    };
    
    onMounted(() => {
      console.log('首页组件已加载');
      
      // 添加滚动事件监听
      window.addEventListener('scroll', handleScroll);
      // 初始检查一次滚动位置
      checkScrollPosition();
      
      // 添加FAQ点击事件
      const faqItems = document.querySelectorAll('.faq-question');
      faqItems.forEach((item, index) => {
        item.addEventListener('click', () => {
          toggleFaq(index);
          const answer = item.nextElementSibling;
          const icon = item.querySelector('i');
          
          if (activeFaq.value === index) {
            answer.style.maxHeight = answer.scrollHeight + 'px';
            icon.classList.add('is-active');
          } else {
            answer.style.maxHeight = '0';
            icon.classList.remove('is-active');
          }
        });
      });
      
      // 添加滚动平滑效果
      document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
          e.preventDefault();
          document.querySelector(this.getAttribute('href')).scrollIntoView({
            behavior: 'smooth'
          });
        });
      });
    });
    
    onUnmounted(() => {
      // 组件卸载时移除事件监听器
      window.removeEventListener('scroll', handleScroll);
      document.body.style.overflow = 'auto'; // 确保离开页面时恢复滚动
    });
    
    return {
      activeFaq,
      activeSection,
      mobileMenuOpen,
      isHeaderVisible,
      toggleFaq,
      toggleMobileMenu
    }
  }
}
</script>

<style lang="scss" scoped>
html, body {
  overflow-y: auto;
  height: auto;
  min-height: 100%;
  margin: 0;
  padding: 0;
  overflow-x: hidden;
  overflow-y: auto !important;
}

.home-container {
  min-height: 100vh;
  background: linear-gradient(125deg, #2a3f5f 0%, #4b6cb7 100%);
  color: white;
  position: relative;
  overflow-x: hidden;
  overflow-y: visible;
  
  &::before, &::after {
    content: '';
    position: absolute;
    border-radius: 50%;
    filter: blur(60px);
    pointer-events: none;
  }
  
  &::before {
    width: 30vw;
    height: 30vw;
    background: rgba(111, 161, 255, 0.25);
    top: -5%;
    right: -5%;
  }
  
  &::after {
    width: 35vw;
    height: 35vw;
    background: rgba(142, 201, 255, 0.25);
    bottom: -10%;
    left: -10%;
  }
}

.background-decoration {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
  
  .bg-shape {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.3;
    pointer-events: none;
  }
  
  .shape-1 {
    width: 500px;
    height: 500px;
    background: linear-gradient(135deg, #7f95d1, transparent);
    top: -200px;
    right: -100px;
    animation: float 30s ease-in-out infinite alternate;
  }
  
  .shape-2 {
    width: 600px;
    height: 600px;
    background: linear-gradient(135deg, #82ccdd, transparent);
    bottom: -300px;
    left: -200px;
    animation: float 35s ease-in-out infinite alternate;
    animation-delay: -5s;
  }
  
  .shape-3 {
    width: 400px;
    height: 400px;
    background: linear-gradient(135deg, #6a89cc, transparent);
    top: 40%;
    right: 30%;
    animation: float 25s ease-in-out infinite alternate;
    animation-delay: -10s;
  }
  
  .bg-dots {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-image: radial-gradient(rgba(255, 255, 255, 0.15) 1px, transparent 1px);
    background-size: 25px 25px;
    opacity: 0.2;
    pointer-events: none;
  }
  
  @keyframes float {
    0% {
      transform: translate(0, 0) rotate(0deg);
    }
    100% {
      transform: translate(30px, 30px) rotate(5deg);
    }
  }
}

.home-content {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 70px 20px 40px;
  min-height: 100vh;
  height: auto;
  display: block;
}

.hero-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 80vh;
  padding: 50px 0;
  overflow: visible;
  
  @media (max-width: 992px) {
    flex-direction: column;
    text-align: center;
    height: auto;
    min-height: 70vh;
  }
  
  .hero-content {
    max-width: 600px;
    animation: fadeIn 0.8s ease-out;
    
    @media (max-width: 992px) {
      max-width: 100%;
      margin-bottom: 50px;
    }
  }
  
  .hero-title {
    font-size: 48px;
    font-weight: 800;
    margin-bottom: 16px;
    line-height: 1.2;
    letter-spacing: -0.5px;
    background: linear-gradient(to right, #ffffff, #82ccdd);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    text-shadow: 0 2px 10px rgba(130, 204, 221, 0.3);
    position: relative;
    display: inline-block;
    
    &::after {
      content: '';
      position: absolute;
      bottom: -8px;
      left: 0;
      width: 60px;
      height: 4px;
      background: linear-gradient(90deg, #6a89cc, #82ccdd);
      border-radius: 2px;
    }
  }
  
  .hero-subtitle {
    font-size: 20px;
    line-height: 1.6;
    color: rgba(255, 255, 255, 0.9);
    margin-bottom: 30px;
    font-weight: 400;
    max-width: 90%;
  }
  
  .hero-actions {
    display: flex;
    gap: 20px;
  }
  
  .hero-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    padding: 16px 32px;
    border-radius: 30px;
    border: none;
    font-size: 17px;
    font-weight: 600;
    color: white;
    transition: all 0.3s ease;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    letter-spacing: 0.8px;
    text-decoration: none;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    text-transform: none;
    
    span {
      position: relative;
      z-index: 2;
      transition: transform 0.3s ease;
      text-shadow: 0 0 1px rgba(255, 255, 255, 0.5);
    }
    
    i {
      font-size: 18px;
      transition: transform 0.3s ease;
      position: relative;
      z-index: 2;
    }
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
      transition: all 0.6s ease;
      z-index: 1;
    }
    
    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
      
      span {
        transform: translateX(-3px);
      }
      
      i {
        transform: translateX(3px);
      }
      
      &::before {
        left: 100%;
      }
    }
    
    &:active {
      transform: translateY(0);
    }
  }
  
  .primary-btn {
    background: linear-gradient(90deg, #6a89cc, #82ccdd);
    
    &:hover {
      box-shadow: 
        0 10px 25px rgba(106, 137, 204, 0.3),
        0 0 20px rgba(130, 204, 221, 0.3);
    }
  }
  
  .secondary-btn {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.3);
    backdrop-filter: blur(5px);
    
    &:hover {
      background: rgba(255, 255, 255, 0.15);
      border-color: rgba(130, 204, 221, 0.5);
      box-shadow: 0 0 15px rgba(130, 204, 221, 0.2);
    }
    
    span {
      color: #ffffff;
      font-weight: 500;
    }
  }
  
  .hero-image {
    width: 40%;
    
    @media (max-width: 992px) {
      width: 80%;
    }
    
    .security-animation {
      width: 100%;
      height: 400px;
      display: flex;
      align-items: center;
      justify-content: center;
      
      .security-shield {
        width: 200px;
        height: 200px;
        background: linear-gradient(135deg, rgba(106, 137, 204, 0.3), rgba(130, 204, 221, 0.3));
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        position: relative;
        box-shadow: 0 0 40px rgba(130, 204, 221, 0.2);
        
        .shield-icon {
          font-size: 60px;
          color: #ffffff;
          animation: pulse 2s infinite;
        }
        
        .shield-pulse {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          border-radius: 50%;
          background: rgba(130, 204, 221, 0.1);
          z-index: -1;
          animation: ripple 2s infinite;
        }
        
        @keyframes pulse {
          0% {
            transform: scale(1);
            opacity: 1;
          }
          50% {
            transform: scale(1.1);
            opacity: 0.8;
          }
          100% {
            transform: scale(1);
            opacity: 1;
          }
        }
        
        @keyframes ripple {
          0% {
            transform: scale(1);
            opacity: 1;
          }
          100% {
            transform: scale(1.5);
            opacity: 0;
          }
        }
      }
    }
  }
}

.features-section {
  padding: 80px 0;
}

.section-title {
  text-align: center;
  font-size: 36px;
  margin-bottom: 20px;
  background: linear-gradient(to right, #6a89cc, #82ccdd);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 1px;
}

.section-subtitle {
  text-align: center;
  font-size: 18px;
  color: #a0aec0;
  margin-bottom: 60px;
  max-width: 700px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.6;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
}

.feature-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  padding: 30px;
  border-radius: 24px;
  color: white;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  animation: fadeIn 1s ease-out;
  animation-fill-mode: both;
  
  &:nth-child(2) {
    animation-delay: 0.2s;
  }
  
  &:nth-child(3) {
    animation-delay: 0.4s;
  }
  
  &:nth-child(4) {
    animation-delay: 0.6s;
  }
  
  &:nth-child(5) {
    animation-delay: 0.8s;
  }
  
  &:nth-child(6) {
    animation-delay: 1s;
  }
  
  &:hover {
    transform: translateY(-5px);
    background: rgba(255, 255, 255, 0.1);
    border-color: rgba(130, 204, 221, 0.4);
    box-shadow: 
      0 15px 30px rgba(0, 0, 0, 0.15),
      0 0 15px rgba(130, 204, 221, 0.2);
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 3px;
    height: 0;
    background: linear-gradient(0deg, #6a89cc, #82ccdd);
    transition: height 0.6s ease;
    filter: blur(1px);
  }
  
  &:hover::before {
    height: 100%;
  }
  
  h3 {
    font-size: 24px;
    margin: 20px 0 15px;
    background: linear-gradient(90deg, #ffffff, #82ccdd);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    font-weight: 700;
  }
  
  p {
    margin: 0 0 20px;
    opacity: 0.9;
    line-height: 1.6;
    font-size: 16px;
    flex-grow: 1;
  }
}

.feature-icon {
  width: 70px;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(130, 204, 221, 0.1);
  border-radius: 20px;
  overflow: hidden;
  position: relative;
  box-shadow: 0 0 20px rgba(130, 204, 221, 0.15);
  
  i {
    font-size: 32px;
    color: #ffffff;
    transition: all 0.3s ease;
  }
  
  &:hover i {
    transform: scale(1.2);
    color: #82ccdd;
  }
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0;
  
  li {
    position: relative;
    padding: 8px 0 8px 15px;
    font-size: 15px;
    opacity: 0.85;
    
    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 15px;
      width: 5px;
      height: 5px;
      border-radius: 50%;
      background: linear-gradient(90deg, #6a89cc, #82ccdd);
    }
    
    &:hover {
      opacity: 1;
      transform: translateX(3px);
      transition: all 0.2s ease;
    }
  }
}

.feature-highlight {
  color: #82ccdd;
  font-weight: 600;
  margin-right: 4px;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    width: 100%;
    height: 1px;
    background: #82ccdd;
    opacity: 0.5;
  }
}

.cta-section {
  padding: 80px 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.cta-content {
  text-align: center;
  max-width: 700px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  padding: 60px;
  border-radius: 30px;
  animation: fadeIn 1s ease-out;
  box-shadow: 
    0 20px 40px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.15) inset,
    0 0 30px rgba(130, 204, 221, 0.15);
    
  h2 {
    font-size: 32px;
    font-weight: 700;
    margin: 0 0 20px;
    background: linear-gradient(to right, #ffffff, #82ccdd);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  
  p {
    font-size: 18px;
    color: #e0e6ed;
    margin: 0 0 30px;
  }
}

.cta-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(50px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 40px;
  }
  
  .hero-subtitle {
    font-size: 18px;
  }
  
  .hero-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .hero-btn {
    width: 100%;
    max-width: 300px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .cta-content {
    padding: 40px 20px;
  }
  
  .cta-actions {
    flex-direction: column;
    align-items: center;
  }
}

.home-header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 100;
  background: rgba(42, 63, 95, 0.75);
  backdrop-filter: blur(15px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  transition: all 0.4s cubic-bezier(0.215, 0.61, 0.355, 1);
  transform: translateY(0);
  
  &.header-hidden {
    transform: translateY(-100%);
    box-shadow: none;
  }
  
  &:hover {
    background: rgba(42, 63, 95, 0.85);
    border-bottom: 1px solid rgba(255, 255, 255, 0.12);
  }
  
  .header-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 15px 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
    position: relative;
    cursor: pointer;
    padding: 5px;
    border-radius: 18px;
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    animation: logoEntrance 1.2s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
    
    &:hover {
      background: rgba(255, 255, 255, 0.05);
      
      .logo-icon {
        transform: rotate(5deg) scale(1.05);
      }
      
      .logo-shield {
        transform: translateY(-2px);
        box-shadow: 0 0 20px rgba(130, 204, 221, 0.6);
      }
      
      .logo-sparkles span {
        animation-duration: 1.5s;
      }
      
      .logo-pulse {
        animation-duration: 1.5s;
      }
      
      .logo-text {
        background-position: right center;
        letter-spacing: 0.8px;
      }
    }
  }
  
  .logo-icon {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: visible;
    transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
  }
  
  .logo-text {
    font-size: 22px;
    font-weight: 700;
    background: linear-gradient(to right, #6a89cc, #82ccdd);
    background-size: 200% auto;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: 0.5px;
    transform: translateY(1px);
    transition: all 0.5s ease;
    text-shadow: 0 0 15px rgba(130, 204, 221, 0.2);
  }
}

.logo-animation {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.logo-shield {
  width: 28px;
  height: 32px;
  background: linear-gradient(135deg, #6a89cc, #82ccdd);
  border-radius: 50% 50% 45% 45% / 60% 60% 40% 40%;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 15px rgba(130, 204, 221, 0.4);
  transform: translateY(0);
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: -10px;
    left: -10px;
    right: -10px;
    bottom: -10px;
    background: linear-gradient(135deg, #6a89cc, #82ccdd);
    opacity: 0.2;
    filter: blur(8px);
    z-index: -1;
  }
  
  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 40%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.4),
      transparent
    );
    transform: skewX(-25deg);
    animation: shieldShine 3s infinite;
  }
}

@keyframes shieldShine {
  0% {
    left: -100%;
  }
  20%, 100% {
    left: 100%;
  }
}

.logo-shield-inner {
  width: 18px;
  height: 18px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  position: relative;
  box-shadow: inset 0 0 5px rgba(106, 137, 204, 0.3);
  
  &::before, &::after {
    content: '';
    position: absolute;
    background: #6a89cc;
    border-radius: 1px;
    box-shadow: 0 0 3px rgba(130, 204, 221, 0.5);
  }
  
  &::before {
    width: 12px;
    height: 2px;
    top: 8px;
    left: 3px;
  }
  
  &::after {
    width: 2px;
    height: 12px;
    top: 3px;
    left: 8px;
  }
}

.logo-pulse {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: rgba(106, 137, 204, 0.2);
  opacity: 0.4;
  transform: scale(1);
  animation: logoPulse 2s infinite cubic-bezier(0.215, 0.61, 0.355, 1);
}

.logo-sparkles {
  position: absolute;
  width: 100%;
  height: 100%;
  
  span {
    position: absolute;
    background: white;
    border-radius: 50%;
    opacity: 0;
    filter: blur(0.5px);
    box-shadow: 0 0 5px rgba(255, 255, 255, 0.8);
    
    &:nth-child(1) {
      width: 4px;
      height: 4px;
      top: -2px;
      right: 0;
      animation: sparkle 2s 0.2s infinite;
    }
    
    &:nth-child(2) {
      width: 3px;
      height: 3px;
      top: 50%;
      right: -3px;
      animation: sparkle 2s 0.5s infinite;
    }
    
    &:nth-child(3) {
      width: 3px;
      height: 3px;
      bottom: -1px;
      left: 30%;
      animation: sparkle 2s 0.8s infinite;
    }
  }
}

@keyframes logoPulse {
  0% {
    opacity: 0.4;
    transform: scale(1);
  }
  50% {
    opacity: 0.2;
    transform: scale(1.3);
  }
  100% {
    opacity: 0;
    transform: scale(1.5);
  }
}

.main-nav {
  .nav-links {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 30px;
    
    .nav-link {
      color: rgba(255, 255, 255, 0.85);
      text-decoration: none;
      font-weight: 500;
      font-size: 16px;
      transition: all 0.4s ease;
      position: relative;
      padding: 5px 2px;
      letter-spacing: 0.2px;
      
      &::after {
        content: '';
        position: absolute;
        bottom: -2px;
        left: 0;
        width: 100%;
        height: 2px;
        background: linear-gradient(90deg, #6a89cc, #82ccdd);
        transform: scaleX(0);
        transform-origin: right;
        transition: transform 0.4s cubic-bezier(0.23, 1, 0.32, 1);
        opacity: 0.8;
      }
      
      &:hover, &.active {
        color: white;
        text-shadow: 0 0 15px rgba(255, 255, 255, 0.4);
        
        &::after {
          transform: scaleX(1);
          transform-origin: left;
        }
      }
      
      &.active {
        font-weight: 600;
        color: #ffffff;
        
        &::after {
          height: 3px;
          box-shadow: 0 0 8px rgba(130, 204, 221, 0.4);
          opacity: 1;
        }
      }
    }
  }
}

.header-actions {
  display: flex;
  gap: 15px;
  
  .header-btn {
    padding: 9px 20px;
    border-radius: 30px;
    font-weight: 600;
    font-size: 15px;
    transition: all 0.3s ease;
    text-decoration: none;
    letter-spacing: 0.3px;
    
    &.login-btn {
      color: white;
      background: linear-gradient(135deg, rgba(106, 137, 204, 0.2), rgba(130, 204, 221, 0.2));
      border: 1px solid rgba(255, 255, 255, 0.2);
      
      &:hover {
        background: linear-gradient(135deg, rgba(106, 137, 204, 0.3), rgba(130, 204, 221, 0.3));
        transform: translateY(-2px);
        border-color: rgba(255, 255, 255, 0.3);
        box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(130, 204, 221, 0.1);
      }
    }
  }
}

.mobile-menu-btn {
  display: none;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  color: white;
  font-size: 22px;
  cursor: pointer;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.15);
    transform: rotate(5deg);
  }
  
  @media (max-width: 768px) {
    display: flex;
  }
}

@media (max-width: 768px) {
  .main-nav, .header-actions {
    display: none;
  }
}

.workflow-section {
  padding: 100px 0;
  
  .workflow-steps {
    display: flex;
    flex-direction: column;
    gap: 30px;
    max-width: 800px;
    margin: 0 auto;
    
    .workflow-step {
      display: flex;
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(10px);
      border-radius: 30px;
      padding: 25px;
      border: 1px solid rgba(255, 255, 255, 0.15);
      box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
      transition: all 0.4s ease;
      
      &:hover {
        transform: translateX(10px);
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(130, 204, 221, 0.4);
        box-shadow: 0 15px 30px rgba(0, 0, 0, 0.15);
      }
      
      .step-number {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        background: linear-gradient(135deg, #6a89cc, #82ccdd);
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 24px;
        font-weight: 700;
        margin-right: 20px;
        flex-shrink: 0;
        box-shadow: 0 5px 15px rgba(106, 137, 204, 0.2);
      }
      
      .step-content {
        flex-grow: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        
        h3 {
          font-size: 22px;
          margin: 0 0 10px;
          background: linear-gradient(90deg, #ffffff, #82ccdd);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
        }
        
        p {
          margin: 0;
          color: #e0e6ed;
          font-size: 16px;
          line-height: 1.5;
        }
        
        .step-icon {
          position: absolute;
          right: 30px;
          font-size: 40px;
          color: rgba(130, 204, 221, 0.2);
          
          @media (max-width: 768px) {
            display: none;
          }
        }
      }
    }
  }
}

.faq-section {
  padding: 100px 0;
  
  .faq-container {
    max-width: 800px;
    margin: 0 auto;
    
    .faq-item {
      margin-bottom: 20px;
      border-radius: 20px;
      overflow: hidden;
      background: rgba(255, 255, 255, 0.05);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.15);
      transition: all 0.3s ease;
      
      &:hover {
        border-color: rgba(130, 204, 221, 0.4);
        box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
      }
      
      .faq-question {
        padding: 20px;
        cursor: pointer;
        display: flex;
        justify-content: space-between;
        align-items: center;
        transition: all 0.3s ease;
        
        &:hover {
          background: rgba(255, 255, 255, 0.1);
        }
        
        h3 {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: white;
        }
        
        i {
          font-size: 20px;
          color: #a0aec0;
          transition: all 0.3s ease;
          
          &.is-active {
            transform: rotate(180deg);
            color: #82ccdd;
          }
        }
      }
      
      .faq-answer {
        padding: 0 20px;
        max-height: 0;
        overflow: hidden;
        transition: max-height 0.3s ease;
        
        p {
          padding-bottom: 20px;
          margin: 0;
          color: #e0e6ed;
          line-height: 1.6;
        }
      }
    }
  }
}

.home-footer {
  padding: 60px 0 30px;
  background: rgba(42, 63, 95, 0.8);
  backdrop-filter: blur(10px);
  margin-top: 50px;
  
  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
  }
  
  .footer-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 40px;
    margin-bottom: 50px;
    
    .footer-col {
      h3 {
        font-size: 18px;
        margin-bottom: 20px;
        font-weight: 600;
        color: white;
      }
    }
  }
  
  .footer-logo {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 15px;
    padding: 5px 8px;
    border-radius: 18px;
    transition: all 0.5s ease;
    position: relative;
    cursor: pointer;
    
    &:hover {
      background: rgba(255, 255, 255, 0.05);
      transform: translateY(-2px);
      
      &::after {
        opacity: 0.8;
        transform: scale(1.05);
      }
      
      .logo-icon {
        transform: rotate(5deg);
        
        .logo-shield {
          box-shadow: 0 0 25px rgba(130, 204, 221, 0.6);
        }
      }
      
      .logo-text {
        text-shadow: 0 0 15px rgba(130, 204, 221, 0.4);
      }
    }
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      border-radius: 18px;
      background: radial-gradient(circle at 30% 30%, rgba(106, 137, 204, 0.1), transparent 70%);
      opacity: 0;
      transition: all 0.5s ease;
      z-index: -1;
    }
    
    .logo-icon {
      width: 36px;
      height: 36px;
      transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
      
      .logo-shield {
        width: 24px;
        height: 28px;
        transition: all 0.5s ease;
      }
      
      .logo-shield-inner {
        width: 16px;
        height: 16px;
        
        &::before {
          width: 10px;
          height: 2px;
          top: 7px;
          left: 3px;
        }
        
        &::after {
          width: 2px;
          height: 10px;
          top: 3px;
          left: 7px;
        }
      }
    }
    
    .logo-text {
      font-size: 18px;
      font-weight: 700;
      background: linear-gradient(to right, #6a89cc, #82ccdd);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      transition: all 0.5s ease;
      letter-spacing: 0.5px;
    }
  }
  
  .footer-desc {
    color: #a0aec0;
    font-size: 14px;
    line-height: 1.6;
    margin-bottom: 20px;
  }
  
  .footer-links {
    list-style: none;
    padding: 0;
    margin: 0;
    
    li {
      margin-bottom: 12px;
      
      a {
        color: #a0aec0;
        text-decoration: none;
        font-size: 14px;
        transition: all 0.3s ease;
        position: relative;
        
        &::after {
          content: '';
          position: absolute;
          bottom: -2px;
          left: 0;
          width: 100%;
          height: 1px;
          background: linear-gradient(90deg, #6a89cc, #82ccdd);
          transform: scaleX(0);
          transform-origin: right;
          transition: transform 0.3s ease;
        }
        
        &:hover {
          color: white;
          
          &::after {
            transform: scaleX(1);
            transform-origin: left;
          }
        }
      }
    }
  }
  
  .footer-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 30px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    
    p {
      color: #a0aec0;
      font-size: 14px;
      margin: 0;
    }
    
    .footer-social {
      display: flex;
      gap: 15px;
      
      .social-link {
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: rgba(255, 255, 255, 0.1);
        border-radius: 50%;
        color: white;
        font-size: 18px;
        transition: all 0.3s ease;
        
        &:hover {
          background: linear-gradient(135deg, #6a89cc, #82ccdd);
          transform: translateY(-3px);
        }
      }
    }
    
    @media (max-width: 768px) {
      flex-direction: column;
      gap: 20px;
      text-align: center;
    }
  }
}

/* 移动端菜单 */
.mobile-menu {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(42, 63, 95, 0.9);
  backdrop-filter: blur(15px);
  z-index: 200;
  opacity: 0;
  visibility: hidden;
  transition: all 0.4s cubic-bezier(0.23, 1, 0.32, 1);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  
  &.mobile-menu-open {
    opacity: 1;
    visibility: visible;
  }
  
  .mobile-menu-content {
    width: 100%;
    max-width: 400px;
    padding: 20px;
    animation: fadeInUp 0.6s ease-out both;
  }
  
  .mobile-menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 40px;
  }
  
  .mobile-menu-close {
    background: rgba(255, 255, 255, 0.05);
    border: none;
    color: white;
    font-size: 24px;
    cursor: pointer;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.3s ease;
    
    &:hover {
      background: rgba(255, 255, 255, 0.15);
      transform: rotate(90deg);
    }
  }
  
  .mobile-nav-links {
    list-style: none;
    padding: 0;
    margin: 0 0 40px;
    
    li {
      margin-bottom: 20px;
      opacity: 0;
      animation: fadeInRight 0.5s ease-out forwards;
      
      &:nth-child(1) { animation-delay: 0.1s; }
      &:nth-child(2) { animation-delay: 0.2s; }
      &:nth-child(3) { animation-delay: 0.3s; }
      
      .mobile-nav-link {
        color: rgba(255, 255, 255, 0.8);
        text-decoration: none;
        font-size: 24px;
        font-weight: 500;
        display: block;
        padding: 10px 0;
        transition: all 0.3s ease;
        position: relative;
        letter-spacing: 0.5px;
        
        &::after {
          content: '';
          position: absolute;
          bottom: 0;
          left: 0;
          width: 30px;
          height: 3px;
          background: linear-gradient(90deg, #6a89cc, #82ccdd);
          transform: scaleX(0);
          transform-origin: left;
          transition: transform 0.4s cubic-bezier(0.23, 1, 0.32, 1), width 0.4s cubic-bezier(0.23, 1, 0.32, 1);
          opacity: 0.8;
        }
        
        &:hover, &.active {
          color: white;
          padding-left: 5px;
          
          &::after {
            transform: scaleX(1);
            width: 60px;
          }
        }
        
        &.active {
          color: white;
          font-weight: 700;
          text-shadow: 0 0 15px rgba(130, 204, 221, 0.5);
        }
      }
    }
  }
  
  .mobile-menu-actions {
    display: flex;
    flex-direction: column;
    gap: 15px;
    opacity: 0;
    animation: fadeInUp 0.5s ease-out forwards;
    animation-delay: 0.4s;
    
    .header-btn {
      width: 100%;
      text-align: center;
      padding: 15px;
      border-radius: 30px;
      font-size: 16px;
      font-weight: 600;
      letter-spacing: 0.5px;
      background: linear-gradient(135deg, rgba(106, 137, 204, 0.2), rgba(130, 204, 221, 0.2));
      color: white;
      border: 1px solid rgba(255, 255, 255, 0.2);
      transition: all 0.3s ease;
      
      &:hover {
        background: linear-gradient(135deg, rgba(106, 137, 204, 0.3), rgba(130, 204, 221, 0.3));
        transform: translateY(-2px);
        border-color: rgba(255, 255, 255, 0.3);
      }
    }
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInRight {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 添加全局样式以确保滚动正常 */
::v-deep #app {
  height: auto !important;
  min-height: 100vh;
  overflow-y: auto !important;
}

/* 确保页面可以滚动 */
::v-deep html, ::v-deep body {
  overflow-y: auto !important;
  height: auto !important;
}

/* 移动端菜单中的LOGO优化 */
.mobile-menu-content {
  .logo {
    padding: 0 8px;
    
    .logo-icon {
      width: 46px;
      height: 46px;
    }
    
    .logo-text {
      font-size: 24px;
      letter-spacing: 0.8px;
    }
    
    .logo-shield {
      width: 32px;
      height: 36px;
    }
    
    .logo-shield-inner {
      width: 22px;
      height: 22px;
      
      &::before {
        width: 14px;
        top: 10px;
        left: 4px;
      }
      
      &::after {
        height: 14px;
        top: 4px;
        left: 10px;
      }
    }
  }
}

/* 页脚LOGO样式优化 */
.footer-logo {
  margin-bottom: 15px;
  
  .logo {
    padding: 5px 8px;
    border-radius: 18px;
    
    .logo-text {
      font-size: 18px;
    }
    
    .logo-icon {
      width: 36px;
      height: 36px;
      
      .logo-shield {
        width: 24px;
        height: 28px;
      }
      
      .logo-shield-inner {
        width: 16px;
        height: 16px;
        
        &::before {
          width: 10px;
          height: 2px;
          top: 7px;
          left: 3px;
        }
        
        &::after {
          width: 2px;
          height: 10px;
          top: 3px;
          left: 7px;
        }
      }
      
      .logo-pulse {
        animation-duration: 3s;
      }
      
      .logo-sparkles span {
        animation-duration: 3s;
        
        &:nth-child(1) {
          animation-delay: 0.4s;
        }
        
        &:nth-child(2) {
          animation-delay: 1s;
        }
        
        &:nth-child(3) {
          animation-delay: 1.6s;
        }
      }
    }
    
    &:hover {
      .logo-shield::after {
        animation-duration: 4s;
      }
    }
  }
}

/* LOGO加载动画 */
@keyframes logoEntrance {
  0% {
    opacity: 0;
    transform: scale(0.8) translateY(-10px);
  }
  40% {
    opacity: 1;
  }
  70% {
    transform: scale(1.05) translateY(3px);
  }
  100% {
    transform: scale(1) translateY(0);
  }
}

/* 盾牌内部'+'符号的脉动效果 */
@keyframes pulseCross {
  0% {
    box-shadow: 0 0 3px rgba(130, 204, 221, 0.5);
  }
  50% {
    box-shadow: 0 0 8px rgba(130, 204, 221, 0.8);
  }
  100% {
    box-shadow: 0 0 3px rgba(130, 204, 221, 0.5);
  }
}

.logo-shield-inner {
  &::before, &::after {
    animation: pulseCross 2s ease-in-out infinite;
  }
}

@keyframes floatShield {
  0% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-2px);
  }
  100% {
    transform: translateY(0);
  }
}

.logo-shield {
  animation: floatShield 3s ease-in-out infinite;
}

@keyframes sparkle {
  0% {
    opacity: 0;
    transform: scale(0);
  }
  50% {
    opacity: 1;
    transform: scale(1);
  }
  100% {
    opacity: 0;
    transform: scale(0);
  }
}

@keyframes shieldShine {
  0% {
    left: -100%;
  }
  20%, 100% {
    left: 100%;
  }
}
</style> 
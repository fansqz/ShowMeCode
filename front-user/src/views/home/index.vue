<template>
  <div class="home">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content">
        <h1 class="hero-title">
          <span class="gradient-text">数据结构</span>
          <br />可视化学习平台
        </h1>
        <p class="hero-subtitle">通过直观的可视化展示，轻松掌握数据结构与算法</p>
        <div class="hero-buttons">
          <el-button type="primary" size="large" class="primary-button" @click="goToCoding">
            开始学习
            <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
          <el-button size="large" class="secondary-button" @click="scrollToFeatures">
            了解更多
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
        </div>
      </div>
      <div class="hero-image">
        <div class="floating-card card-1">
          <el-icon><DataLine /></el-icon>
          <span>可视化学习</span>
        </div>
        <div class="floating-card card-2">
          <el-icon><Monitor /></el-icon>
          <span>在线调试</span>
        </div>
        <div class="floating-card card-3">
          <el-icon><Reading /></el-icon>
          <span>知识库</span>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features" ref="featuresSection">
      <div class="section-header">
        <h2 class="section-title">核心功能</h2>
        <p class="section-subtitle">探索我们的特色功能，开启你的编程学习之旅</p>
      </div>
      <div class="features-container">
        <div class="feature-item" v-for="(feature, index) in features" :key="index">
          <div class="feature-icon-wrapper">
            <el-icon class="feature-icon" :size="32">
              <component :is="feature.icon" />
            </el-icon>
          </div>
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.description }}</p>
        </div>
      </div>
    </section>

    <!-- Showcase Section -->
    <section class="showcase">
      <div class="section-header">
        <h2 class="section-title">功能展示</h2>
        <p class="section-subtitle">直观感受我们的平台特色</p>
      </div>
      <div class="showcase-container">
        <div class="showcase-item" v-for="(item, index) in showcaseItems" :key="index">
          <div class="showcase-image">
            <img :src="item.image" :alt="item.title" />
          </div>
          <div class="showcase-content">
            <h3>{{ item.title }}</h3>
            <p>{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta">
      <div class="cta-content">
        <h2>准备好开始你的编程之旅了吗？</h2>
        <p>立即加入我们，开启数据结构与算法的学习之旅</p>
        <el-button type="primary" size="large" class="cta-button" @click="goToCoding">
          立即开始
          <el-icon class="el-icon--right"><ArrowRight /></el-icon>
        </el-button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { DataLine, Monitor, Reading, SetUp, ArrowRight, ArrowDown } from '@element-plus/icons-vue'

  const router = useRouter()
  const featuresSection = ref<HTMLElement | null>(null)

  const features = ref([
    {
      icon: DataLine,
      title: '数据结构可视化',
      description: '通过动态演示，直观理解各种数据结构的实现原理和操作过程',
    },
    {
      icon: Monitor,
      title: '在线调试',
      description: '提供实时编译和运行环境，支持断点调试，快速定位问题',
    },
    {
      icon: Reading,
      title: '知识库学习',
      description: '系统化的学习资料，从基础到进阶，循序渐进掌握编程知识',
    },
    {
      icon: SetUp,
      title: '自定义数据结构',
      description: '支持创建和测试自定义数据结构，培养编程思维和创新能力',
    },
  ])

  const showcaseItems = ref([
    {
      title: '可视化学习体验',
      description: '通过直观的图形界面，轻松理解复杂的数据结构概念',
      image: '/images/showcase-visualization.png',
    },
    {
      title: '实时代码调试',
      description: '在线编写和调试代码，即时查看运行结果和数据结构变化',
      image: '/images/showcase-debug.png',
    },
  ])

  const goToCoding = () => {
    router.push('/coding')
  }

  const scrollToFeatures = () => {
    featuresSection.value?.scrollIntoView({ behavior: 'smooth' })
  }
</script>

<style scoped lang="scss">
  .home {
    width: 100%;
    min-height: 100vh;
    overflow-x: hidden;
    background: linear-gradient(to bottom, $bg-color-secondary, $base-background-color);
    transition: all 0.3s ease;
  }

  .hero {
    position: relative;
    box-sizing: border-box;
    display: flex;
    gap: 40px;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    max-width: 1200px;
    padding: 80px 20px;
    margin: 0 auto;
    overflow: hidden;

    .hero-content {
      flex: 1;
      max-width: 600px;
    }

    .hero-title {
      margin-bottom: 20px;
      font-size: 4em;
      font-weight: 800;
      line-height: 1.2;
      color: $text-color-primary;
      transition: color 0.3s ease;

      .gradient-text {
        background: linear-gradient(45deg, $primary-color, #36d1dc);
        background-clip: text;
        -webkit-text-fill-color: transparent;
      }
    }

    .hero-subtitle {
      margin-bottom: 40px;
      font-size: 1.5em;
      line-height: 1.6;
      color: $text-color-regular;
    }

    .hero-buttons {
      display: flex;
      gap: 20px;
    }

    .primary-button {
      padding: 15px 40px;
      font-size: 1.2em;
      border-radius: 30px;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 4px 15px rgb(64 158 255 / 30%);
        transform: translateY(-2px);
      }
    }

    .secondary-button {
      padding: 15px 40px;
      font-size: 1.2em;
      border-radius: 30px;
      transition: all 0.3s ease;

      &:hover {
        background: $fill-color-light;
        transform: translateY(-2px);
      }
    }

    .hero-image {
      position: relative;
      flex: 1;
      max-width: 600px;
      overflow: visible;

      img {
        display: block;
        width: 100%;
        height: auto;
        border-radius: 20px;
        box-shadow: 0 20px 40px rgb(0 0 0 / 10%);
      }

      .floating-card {
        position: absolute;
        z-index: 1;
        display: flex;
        gap: 10px;
        align-items: center;
        padding: 15px 25px;
        background: $base-background-color;
        border: 1px solid $border-color-light;
        border-radius: 15px;
        box-shadow: $box-shadow;
        animation: float 3s ease-in-out infinite;

        .el-icon {
          font-size: 20px;
          color: $primary-color;
        }

        span {
          font-weight: 500;
          color: $text-color-primary;
        }
      }

      .card-1 {
        top: 20%;
        left: -10px;
        animation-delay: 0s;
      }

      .card-2 {
        top: 50%;
        right: -10px;
        animation-delay: 0.5s;
      }

      .card-3 {
        bottom: 20%;
        left: 50%;
        transform: translateX(-50%);
        animation-delay: 1s;
      }
    }
  }

  .section-header {
    margin-bottom: 60px;
    text-align: center;

    .section-title {
      margin-bottom: 20px;
      font-size: 2.5em;
      font-weight: 700;
      color: $text-color-primary;
    }

    .section-subtitle {
      max-width: 600px;
      margin: 0 auto;
      font-size: 1.2em;
      color: $text-color-regular;
    }
  }

  .features {
    box-sizing: border-box;
    width: 100%;
    padding: 100px 0;
    background: $base-background-color;

    .features-container {
      display: flex;
      gap: 30px;
      justify-content: space-between;
      max-width: 1200px;
      padding: 0 20px;
      margin: 0 auto;
    }

    .feature-item {
      flex: 1;
      padding: 40px 30px;
      text-align: center;
      background: $base-background-color;
      border: 1px solid $border-color-light;
      border-radius: $border-radius-xl;
      box-shadow: $box-shadow-light;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: $box-shadow;
        transform: translateY(-5px);
      }

      .feature-icon-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 80px;
        height: 80px;
        margin: 0 auto 25px;
        background: $fill-color-light;
        border-radius: 50%;
      }

      .feature-icon {
        color: $primary-color;
      }

      h3 {
        margin-bottom: 15px;
        font-size: 1.4em;
        font-weight: 600;
        color: $text-color-primary;
      }

      p {
        font-size: 1.1em;
        line-height: 1.6;
        color: $text-color-regular;
      }
    }
  }

  .showcase {
    box-sizing: border-box;
    width: 100%;
    padding: 100px 0;
    background: $bg-color-secondary;

    .showcase-container {
      max-width: 1200px;
      padding: 0 20px;
      margin: 0 auto;
    }

    .showcase-item {
      display: flex;
      align-items: center;
      margin-bottom: 80px;
      overflow: hidden;
      background: $base-background-color;
      border: 1px solid $border-color-light;
      border-radius: 20px;
      box-shadow: $box-shadow-light;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: $box-shadow;
        transform: translateY(-5px);
      }

      &:nth-child(even) {
        flex-direction: row-reverse;
      }

      .showcase-image {
        flex: 1;
        height: 400px;
        overflow: hidden;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.3s ease;
        }

        &:hover img {
          transform: scale(1.05);
        }
      }

      .showcase-content {
        flex: 1;
        padding: 40px;
        text-align: left;

        h3 {
          margin-bottom: 20px;
          font-size: 1.8em;
          font-weight: 600;
          color: $text-color-primary;
        }

        p {
          font-size: 1.1em;
          line-height: 1.6;
          color: $text-color-regular;
        }
      }
    }
  }

  .cta {
    padding: 100px 0;
    color: white;
    text-align: center;
    background: linear-gradient(45deg, $primary-color, #36d1dc);

    .cta-content {
      max-width: 800px;
      padding: 0 20px;
      margin: 0 auto;

      h2 {
        margin-bottom: 20px;
        font-size: 2.5em;
        font-weight: 700;
      }

      p {
        margin-bottom: 40px;
        font-size: 1.2em;
        opacity: 0.9;
      }

      .cta-button {
        padding: 15px 40px;
        font-size: 1.2em;
        color: $primary-color;
        background: $base-background-color;
        border: none;
        border-radius: 30px;
        transition: all 0.3s ease;

        &:hover {
          box-shadow: 0 4px 15px rgb(255 255 255 / 30%);
          transform: translateY(-2px);
        }
      }
    }
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0);
    }

    50% {
      transform: translateY(-10px);
    }
  }

  @media (width <= 768px) {
    .hero {
      flex-direction: column;
      padding: 40px 20px;
      overflow: visible;
      text-align: center;

      .hero-title {
        font-size: 3em;
      }

      .hero-buttons {
        justify-content: center;
      }

      .hero-image {
        margin-top: 40px;
        overflow: visible;
      }
    }

    .showcase-item {
      flex-direction: column !important;

      .showcase-image {
        height: 300px;
      }

      .showcase-content {
        padding: 30px;
        text-align: center;
      }
    }

    .features-container {
      flex-direction: column;
    }

    .feature-item {
      margin-bottom: 30px;
    }
  }
</style>

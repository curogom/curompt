import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';
import Translate from '@docusaurus/Translate';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">
          <Translate id="homepage.hero.line1">
            LLM 프롬프트를 분석 · 평가 · 최적화하는 로컬‑퍼스트 CLI
          </Translate>
          <br />
          <Translate id="homepage.hero.line2">
            네트워크로 프롬프트 내용을 전송하지 않는 안전한 로컬 기반 파이프라인
          </Translate>
        </p>
        <div className={styles.buttons}>
          <Link className="button button--secondary button--lg" to="/docs/intro">
            <Translate id="homepage.hero.cta.docs">문서 보기</Translate>
          </Link>
          <Link
            className="button button--outline button--lg"
            to="https://github.com/curogom/curompt">
            <Translate id="homepage.hero.cta.github">GitHub</Translate>
          </Link>
        </div>
        <pre className={styles.heroSnippet}>
          <code>brew install curompt</code>
        </pre>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={siteConfig.title}
      description="curompt: analyze, evaluate, and optimize LLM prompts from the CLI">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}

import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import Translate, {translate} from '@docusaurus/Translate';

type FeatureItem = {
  title: string;
  image: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: translate({
      id: 'homepage.features.cli.title',
      message: 'CLI 기록 기반 분석',
    }),
    image: require('@site/static/img/local_analysis.png').default,
    description: (
      <Translate id="homepage.features.cli.desc">
        터미널에서 수집한 프롬프트 실행 기록을 분석해 섹션 누락이나 금지어 위반을
        미리 발견할 수 있습니다.
      </Translate>
    ),
  },
  {
    title: translate({
      id: 'homepage.features.eval.title',
      message: '멀티샘플 평가 자동화',
    }),
    image: require('@site/static/img/multisample_eval.png').default,
    description: (
      <Translate id="homepage.features.eval.desc">
        Claude·OpenAI·Gemini 등 공급자별로 다중 샘플을 수집해 스키마 적합률, 지연,
        비용을 자동으로 기록합니다.
      </Translate>
    ),
  },
  {
    title: translate({
      id: 'homepage.features.report.title',
      message: '증거 기반 리포트',
    }),
    image: require('@site/static/img/reporting.png').default,
    description: (
      <Translate id="homepage.features.report.desc">
        점수·토큰 사용량·개선 제안을 Markdown/JSON 리포트로 생성해 PR이나 문서에
        근거로 첨부할 수 있습니다.
      </Translate>
    ),
  },
];

function Feature({title, image, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <img className={styles.featureImg} src={image} alt={title} />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

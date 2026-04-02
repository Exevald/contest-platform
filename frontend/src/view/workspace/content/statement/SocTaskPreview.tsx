import {Badge} from '../../../../common/components/badge/Badge'
import {Card} from '../../../../common/components/card/Card'
import styles from './TaskPreview.module.css'

type SocTaskPreviewProps = {
	taskId: string,
}

function SocTaskPreview({taskId}: SocTaskPreviewProps) {
	switch (taskId) {
		case 'soc-task-1':
			return <SocTaskOnePreview />
		case 'soc-task-2':
			return <SocTaskTwoPreview />
		case 'soc-task-3':
			return <SocTaskThreePreview />
		default:
			return null
	}
}

function SocTaskOnePreview() {
	const logs = [
		{severity: 'CRITICAL', tone: 'error' as const, payload: 'sourceIp=10.0.0.1, userId=db_user, eventType=SQL_INJECTION'},
		{severity: 'HIGH', tone: 'warning' as const, payload: 'sourceIp=192.168.1.50, userId=alice, eventType=Brute Force'},
		{severity: 'MEDIUM', tone: 'neutral' as const, payload: 'sourceIp=172.16.0.8, userId=ops, eventType=failed login'},
	]

	return (
		<div className={styles.previewStack}>
			<div className={styles.previewHeroSoc}>
				<div>
					<div className={styles.previewEyebrowSoc}>SOC OUTPUT PREVIEW</div>
					<div className={styles.previewTitleSoc}>Unified Log Stream</div>
				</div>
				<Badge tone="error">Integrity 100%</Badge>
			</div>

			<Card className={styles.socStreamCard}>
				<div className={styles.socMatrix}>
					<div className={styles.socMetric}>
						<span>UNKNOWN</span>
						<strong>0</strong>
					</div>
					<div className={styles.socMetric}>
						<span>Deduplicated</span>
						<strong>12</strong>
					</div>
					<div className={styles.socMetric}>
						<span>UTC Normalized</span>
						<strong>24</strong>
					</div>
				</div>
				<div className={styles.logList}>
					{logs.map(log => (
						<div key={log.payload} className={styles.logRow}>
							<Badge tone={log.tone}>{log.severity}</Badge>
							<div className={styles.logPayload}>{log.payload}</div>
							<div className={styles.logPulse} />
						</div>
					))}
				</div>
			</Card>
		</div>
	)
}

function SocTaskTwoPreview() {
	const suspicious = [
		'85.214.132.117',
		'177.43.10.9',
		'10.15.17.200',
	]

	return (
		<div className={styles.previewStack}>
			<div className={styles.previewHeroSoc}>
				<div>
					<div className={styles.previewEyebrowSoc}>SOC OUTPUT PREVIEW</div>
					<div className={styles.previewTitleSoc}>Suspicious IP Radar</div>
				</div>
				<Badge tone="warning">Score &gt; 50</Badge>
			</div>

			<div className={styles.radarGrid}>
				<Card className={styles.radarCard}>
					<div className={styles.radarStage}>
						<div className={styles.ringOne} />
						<div className={styles.ringTwo} />
						<div className={styles.ringThree} />
						<div className={styles.nodeCore} />
						<div className={styles.nodeDangerA}>85.214</div>
						<div className={styles.nodeDangerB}>177.43</div>
						<div className={styles.nodeDangerC}>10.15</div>
						<div className={styles.nodeSafeA} />
						<div className={styles.nodeSafeB} />
					</div>
				</Card>

				<Card className={styles.suspiciousListCard}>
					<div className={styles.dispatchLabel}>Detected</div>
					{suspicious.map(ip => (
						<div key={ip} className={styles.suspiciousRow}>
							<span>{ip}</span>
							<Badge tone="error">BRUTE_FORCE</Badge>
						</div>
					))}
				</Card>
			</div>
		</div>
	)
}

function SocTaskThreePreview() {
	return (
		<div className={styles.previewStack}>
			<div className={styles.previewHeroSoc}>
				<div>
					<div className={styles.previewEyebrowSoc}>SOC OUTPUT PREVIEW</div>
					<div className={styles.previewTitleSoc}>WAF Mitigation Dashboard</div>
				</div>
				<Badge tone="ok">SYSTEM STATUS: SECURED</Badge>
			</div>

			<div className={styles.metricsGrid}>
				<Card className={styles.metricCardSoc}>
					<div className={styles.metricLabelSoc}>Mitigation Rate</div>
					<div className={styles.metricValueSoc}>96.40%</div>
					<div className={styles.metricTrack}>
						<div className={styles.metricFillGood} style={{width: '96.4%'}} />
					</div>
				</Card>
				<Card className={styles.metricCardSoc}>
					<div className={styles.metricLabelSoc}>False Positive Rate</div>
					<div className={styles.metricValueSoc}>0.42%</div>
					<div className={styles.metricTrack}>
						<div className={styles.metricFillLow} style={{width: '8%'}} />
					</div>
				</Card>
			</div>

			<div className={styles.metricsGridWide}>
				<Card className={styles.rulesCard}>
					<div className={styles.dispatchLabel}>Generated Rules</div>
					<pre className={styles.codeSoc}>
						<code>{`RULE_001: BLOCK\nTARGET: IP_RANGE\nCONDITION: 0\nACTION: DROP\n\nRULE_002: RATE_LIMIT\nTARGET: USER_AGENT\nCONDITION: 10\nACTION: CHALLENGE`}</code>
					</pre>
				</Card>
				<Card className={styles.rulesCard}>
					<div className={styles.dispatchLabel}>Traffic Simulation</div>
					<div className={styles.simRow}>
						<span>Blocked Attacks</span>
						<div className={styles.metricTrack}><div className={styles.metricFillGood} style={{width: '96%'}} /></div>
					</div>
					<div className={styles.simRow}>
						<span>Legitimate Dropped</span>
						<div className={styles.metricTrack}><div className={styles.metricFillLow} style={{width: '2%'}} /></div>
					</div>
				</Card>
			</div>
		</div>
	)
}

export {
	SocTaskPreview,
}

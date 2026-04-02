import styles from './Preloader.module.css'

const Preloader = ({
                       size = 40,
                       strokeWidth = 4,
                       speed = 1,
                       fullScreen = true,
                   }) => {
    const radius = 20
    const circumference = 2 * Math.PI * radius
    const arcLength = circumference / 2

    const loaderElement = (
        <svg
            className={styles.loader}
            viewBox="0 0 50 50"
            style={{
                width: size,
                height: size,
                animationDuration: `${speed}s`,
            }}
        >
            <circle
                className={styles.arc}
                cx="25"
                cy="25"
                r={radius}
                strokeWidth={strokeWidth}
                strokeDasharray={arcLength}
            />
        </svg>
    )

    if (fullScreen) {
        return <div className={styles.wrapper}>{loaderElement}</div>
    }

    return loaderElement
}

export {Preloader}

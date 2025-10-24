'use client';

import { BallSummary } from '@/store/reducers/scorecardSlice';

interface HorizontalBallDisplayProps {
    ball: BallSummary;
    overNumber: number;
    ballNumber: number;
    currentScore: string;
    currentWickets: number;
}

export function HorizontalBallDisplay({
    ball,
    overNumber,
    ballNumber,
    currentScore,
    currentWickets,
}: HorizontalBallDisplayProps): React.JSX.Element {
    // Format the over display (e.g., "4.1" for over 4, ball 1)
    const overDisplay = `${overNumber}.${ballNumber}`;

    // Format the score display (e.g., "47/1")
    const scoreDisplay = `${currentScore}/${currentWickets}`;

    // Determine the ball event display based on ball type and runs
    const getBallEventDisplay = (): string => {
        const isWicket = ball.is_wicket;

        if (isWicket) {
            return 'W';
        }

        // Handle byes with different ball types
        if (ball.byes > 0) {
            const byeText = `B${ball.byes}`;

            if (ball.ball_type === 'NO_BALL' || ball.ball_type === 'no_ball') {
                return `${byeText} + nb`;
            }
            if (ball.ball_type === 'WIDE' || ball.ball_type === 'wide') {
                return `${byeText} + wd`;
            }
            if (ball.run_type === 'LB') {
                return `${byeText} + lb`;
            }
            // Regular byes on good ball
            return byeText;
        }

        // Handle other ball types without byes
        switch (ball.ball_type) {
            case 'WIDE':
            case 'wide':
                return 'Wd';
            case 'NO_BALL':
            case 'no_ball':
                return 'Nb';
            case 'DEAD_BALL':
            case 'dead_ball':
                return 'Db';
            default:
                // Handle run types
                switch (ball.run_type) {
                    case 'LB':
                        return `Lb + ${ball.runs || 0}`;
                    case 'WC':
                        return 'W';
                    default:
                        return ball.runs?.toString() || '0';
                }
        }
    };

    const ballEventDisplay = getBallEventDisplay();

    return (
        <div className="flex items-center justify-between w-full p-3 border-b border-gray-200 hover:bg-gray-50 transition-colors">
            <div className="flex items-center space-x-6 w-full">
                {/* Ball Event */}
                <div className="w-20 text-center">
                    <span className="text-sm font-semibold text-gray-900 bg-gray-100 px-2 py-1 rounded">
                        {ballEventDisplay}
                    </span>
                </div>

                {/* Over */}
                <div className="w-16 text-center">
                    <span className="text-sm font-medium text-gray-700">
                        {overDisplay}
                    </span>
                </div>

                {/* Score */}
                <div className="w-20 text-center">
                    <span className="text-sm font-medium text-gray-700 bg-blue-50 px-2 py-1 rounded">
                        {scoreDisplay}
                    </span>
                </div>
            </div>
        </div>
    );
}

interface HorizontalBallListProps {
    balls: BallSummary[];
    overNumber: number;
    currentScore: number;
    currentWickets: number;
}

export function HorizontalBallList({
    balls,
    overNumber,
    currentScore,
    currentWickets,
}: HorizontalBallListProps): React.JSX.Element {
    // Sort balls by ball number
    const sortedBalls = [...balls].sort((a, b) => a.ball_number - b.ball_number);

    // Track running score and wickets
    let runningScore = 0;
    let runningWickets = 0;

    return (
        <div className="w-full">
            <div className="bg-gray-100 p-3 border-b border-gray-300">
                <div className="flex items-center justify-between text-sm font-semibold text-gray-700">
                    <div className="w-20 text-center">Ball</div>
                    <div className="w-16 text-center">Over</div>
                    <div className="w-20 text-center">Score</div>
                </div>
            </div>

            <div className="max-h-64 overflow-y-auto">
                {sortedBalls.map((ball, index) => {
                    // Update running totals
                    runningScore += (ball.runs || 0) + (ball.byes || 0);
                    if (ball.is_wicket) {
                        runningWickets += 1;
                    }

                    return (
                        <HorizontalBallDisplay
                            key={`${ball.ball_number}-${index}`}
                            ball={ball}
                            overNumber={overNumber}
                            ballNumber={ball.ball_number}
                            currentScore={runningScore.toString()}
                            currentWickets={runningWickets}
                        />
                    );
                })}
            </div>
        </div>
    );
}

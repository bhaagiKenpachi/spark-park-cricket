'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchScorecardRequest,
    startScoringRequest,
    addBallRequest,
    clearScorecard,
    ScorecardResponse,
    BallEventRequest,
    BallType,
    RunType,
    WicketType
} from '@/store/reducers/scorecardSlice';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ArrowLeft, Play, Plus } from 'lucide-react';

interface ScorecardViewProps {
    matchId: string;
    onBack: () => void;
}

export function ScorecardView({ matchId, onBack }: ScorecardViewProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { scorecard, loading, error, scoring } = useAppSelector((state) => state.scorecard);

    const [showAddBallForm, setShowAddBallForm] = useState(false);
    const [ballForm, setBallForm] = useState<Partial<BallEventRequest>>({
        match_id: matchId,
        innings_number: 1,
        ball_type: 'good',
        run_type: '0',
        is_wicket: false,
        byes: 0,
    });

    useEffect(() => {
        dispatch(fetchScorecardRequest(matchId));
        return () => {
            dispatch(clearScorecard());
        };
    }, [dispatch, matchId]);

    const handleStartScoring = () => {
        dispatch(startScoringRequest(matchId));
    };

    const handleAddBall = () => {
        if (ballForm.match_id && ballForm.innings_number && ballForm.ball_type && ballForm.run_type !== undefined) {
            dispatch(addBallRequest(ballForm as BallEventRequest));
            setShowAddBallForm(false);
            // Reset form
            setBallForm({
                match_id: matchId,
                innings_number: ballForm.innings_number,
                ball_type: 'good',
                run_type: '0',
                is_wicket: false,
                byes: 0,
            });
        }
    };

    const getRunTypeOptions = (): { value: RunType; label: string }[] => [
        { value: '0', label: '0 (Dot Ball)' },
        { value: '1', label: '1 Run' },
        { value: '2', label: '2 Runs' },
        { value: '3', label: '3 Runs' },
        { value: '4', label: '4 Runs' },
        { value: '5', label: '5 Runs' },
        { value: '6', label: '6 Runs' },
        { value: 'NB', label: 'No Ball' },
        { value: 'WD', label: 'Wide' },
        { value: 'LB', label: 'Leg Byes' },
        { value: 'WC', label: 'Wicket' },
    ];

    const getBallTypeOptions = (): { value: BallType; label: string }[] => [
        { value: 'good', label: 'Good Ball' },
        { value: 'wide', label: 'Wide' },
        { value: 'no_ball', label: 'No Ball' },
        { value: 'dead_ball', label: 'Dead Ball' },
    ];

    const getWicketTypeOptions = (): { value: WicketType; label: string }[] => [
        { value: 'bowled', label: 'Bowled' },
        { value: 'caught', label: 'Caught' },
        { value: 'lbw', label: 'LBW' },
        { value: 'run_out', label: 'Run Out' },
        { value: 'stumped', label: 'Stumped' },
        { value: 'hit_wicket', label: 'Hit Wicket' },
    ];


    if (loading && !scorecard) {
        return (
            <div className="w-full max-w-6xl mx-auto p-6">
                <div className="flex items-center justify-center py-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                    <span className="ml-2 text-sm text-gray-600">Loading scorecard...</span>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="w-full max-w-6xl mx-auto p-6">
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
                    <strong className="font-bold">Error:</strong>
                    <span className="block sm:inline"> {error}</span>
                    <div className="mt-3">
                        <Button
                            onClick={() => dispatch(fetchScorecardRequest(matchId))}
                            variant="destructive"
                            size="sm"
                        >
                            Retry
                        </Button>
                    </div>
                </div>
            </div>
        );
    }

    if (!scorecard || !scorecard.data) {
        return (
            <div className="w-full max-w-6xl mx-auto p-6">
                <div className="text-center py-8">
                    <p className="text-muted-foreground mb-4">No scorecard found for this match.</p>
                    <Button onClick={onBack}>Back to Matches</Button>
                </div>
            </div>
        );
    }

    // Extract the actual scorecard data from the nested structure
    const scorecardData = scorecard.data;

    return (
        <div className="w-full max-w-6xl mx-auto p-6">
            {/* Header */}
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6 gap-4">
                <div className="flex items-center space-x-4">
                    <Button variant="outline" onClick={onBack}>
                        <ArrowLeft className="h-4 w-4 mr-2" />
                        Back
                    </Button>
                    <div>
                        <h1 className="text-2xl lg:text-3xl font-bold">
                            {scorecardData.series_name} - Match #{scorecardData.match_number}
                        </h1>
                        <p className="text-sm lg:text-base text-gray-600">
                            {scorecardData.team_a} vs {scorecardData.team_b}
                        </p>
                    </div>
                </div>
                <div className="flex space-x-2">
                    {scorecardData.match_status === 'live' && (
                        <Button
                            onClick={() => setShowAddBallForm(true)}
                            disabled={scoring}
                            size="lg"
                        >
                            <Plus className="h-4 w-4 mr-2" />
                            Add Ball
                        </Button>
                    )}
                    {scorecardData.match_status !== 'live' && (
                        <Button
                            onClick={handleStartScoring}
                            disabled={scoring}
                            size="lg"
                        >
                            <Play className="h-4 w-4 mr-2" />
                            Start Scoring
                        </Button>
                    )}
                </div>
            </div>

            {/* Match Info */}
            <Card className="mb-6">
                <CardHeader>
                    <CardTitle>Match Information</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 text-sm">
                        <div className="flex flex-col">
                            <span className="font-medium text-gray-600 mb-1">Total Overs</span>
                            <span className="text-lg font-semibold">{scorecardData.total_overs}</span>
                        </div>
                        <div className="flex flex-col">
                            <span className="font-medium text-gray-600 mb-1">Toss</span>
                            <span className="text-lg font-semibold">Team {scorecardData.toss_winner} ({scorecardData.toss_type === 'H' ? 'Heads' : 'Tails'})</span>
                        </div>
                        <div className="flex flex-col">
                            <span className="font-medium text-gray-600 mb-1">Current Innings</span>
                            <span className="text-lg font-semibold">{scorecardData.current_innings}</span>
                        </div>
                        <div className="flex flex-col">
                            <span className="font-medium text-gray-600 mb-1">Status</span>
                            <Badge
                                variant={scorecardData.match_status === 'live' ? 'default' : 'secondary'}
                                className="w-fit"
                            >
                                {scorecardData.match_status}
                            </Badge>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* Innings */}
            {scorecardData.innings && Array.isArray(scorecardData.innings) && scorecardData.innings.length > 0 ? (
                scorecardData.innings.map((innings) => (
                    <Card key={innings.innings_number} className="mb-6">
                        <CardHeader>
                            <CardTitle>
                                Innings {innings.innings_number} - Team {innings.batting_team}
                                <Badge
                                    variant={innings.status === 'completed' ? 'secondary' : 'default'}
                                    className="ml-2"
                                >
                                    {innings.status}
                                </Badge>
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            {/* Innings Summary */}
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-6 mb-6">
                                <div className="flex flex-col">
                                    <span className="font-medium text-gray-600 mb-1">Runs</span>
                                    <span className="text-2xl font-bold text-green-600">{innings.total_runs}</span>
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-medium text-gray-600 mb-1">Wickets</span>
                                    <span className="text-2xl font-bold text-red-600">{innings.total_wickets}</span>
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-medium text-gray-600 mb-1">Overs</span>
                                    <span className="text-2xl font-bold text-blue-600">{innings.total_overs}</span>
                                </div>
                                <div className="flex flex-col">
                                    <span className="font-medium text-gray-600 mb-1">Balls</span>
                                    <span className="text-2xl font-bold text-purple-600">{innings.total_balls}</span>
                                </div>
                            </div>

                            {/* Extras and Overs in horizontal layout */}
                            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                                {/* Extras */}
                                {innings.extras && (
                                    <div>
                                        <h4 className="font-medium mb-3 text-gray-700">Extras</h4>
                                        <div className="grid grid-cols-2 gap-3">
                                            <div className="flex justify-between">
                                                <span className="text-sm text-gray-600">Byes:</span>
                                                <span className="font-semibold">{innings.extras.byes}</span>
                                            </div>
                                            <div className="flex justify-between">
                                                <span className="text-sm text-gray-600">Leg Byes:</span>
                                                <span className="font-semibold">{innings.extras.leg_byes}</span>
                                            </div>
                                            <div className="flex justify-between">
                                                <span className="text-sm text-gray-600">Wides:</span>
                                                <span className="font-semibold">{innings.extras.wides}</span>
                                            </div>
                                            <div className="flex justify-between">
                                                <span className="text-sm text-gray-600">No Balls:</span>
                                                <span className="font-semibold">{innings.extras.no_balls}</span>
                                            </div>
                                            <div className="flex justify-between col-span-2 pt-2 border-t">
                                                <span className="font-medium text-gray-700">Total:</span>
                                                <span className="font-bold text-lg">{innings.extras.total}</span>
                                            </div>
                                        </div>
                                    </div>
                                )}

                                {/* Overs */}
                                <div>
                                    <h4 className="font-medium mb-3 text-gray-700">Overs</h4>
                                    <div className="space-y-2">
                                        {innings.overs && Array.isArray(innings.overs) && innings.overs.length > 0 ? (
                                            innings.overs.map((over) => (
                                                <div key={over.over_number} className="border rounded p-3 bg-gray-50">
                                                    <div className="flex items-center justify-between mb-2">
                                                        <span className="font-medium">Over {over.over_number}</span>
                                                        <div className="text-sm text-gray-600">
                                                            {over.total_runs} runs, {over.total_wickets} wickets
                                                        </div>
                                                    </div>
                                                    <div className="flex flex-wrap gap-1">
                                                        {over.balls && Array.isArray(over.balls) ? over.balls.map((ball) => (
                                                            <Badge
                                                                key={ball.ball_number}
                                                                variant={ball.is_wicket ? 'destructive' : 'outline'}
                                                                className="text-xs"
                                                            >
                                                                {ball.run_type}
                                                                {ball.byes > 0 && `+${ball.byes}`}
                                                            </Badge>
                                                        )) : (
                                                            <span className="text-xs text-gray-500">No balls</span>
                                                        )}
                                                    </div>
                                                </div>
                                            ))
                                        ) : (
                                            <div className="text-center text-gray-500 py-4">
                                                <p className="text-sm">No overs played yet.</p>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                ))
            ) : (
                <Card className="mb-6">
                    <CardContent className="pt-6">
                        <div className="text-center text-gray-500">
                            <p>No innings data available yet.</p>
                            <p className="text-sm mt-1">Start scoring to see innings details.</p>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Add Ball Form */}
            {showAddBallForm && (
                <Card className="fixed inset-4 bg-white z-50 overflow-auto">
                    <CardHeader>
                        <CardTitle>Add Ball</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <Label htmlFor="innings">Innings</Label>
                                <Select
                                    value={ballForm.innings_number?.toString()}
                                    onValueChange={(value) => setBallForm(prev => ({ ...prev, innings_number: parseInt(value) }))}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="1">1</SelectItem>
                                        <SelectItem value="2">2</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>

                            <div>
                                <Label htmlFor="ball_type">Ball Type</Label>
                                <Select
                                    value={ballForm.ball_type}
                                    onValueChange={(value) => setBallForm(prev => ({ ...prev, ball_type: value as BallType }))}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getBallTypeOptions().map(option => (
                                            <SelectItem key={option.value} value={option.value}>
                                                {option.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            <div>
                                <Label htmlFor="run_type">Run Type</Label>
                                <Select
                                    value={ballForm.run_type}
                                    onValueChange={(value) => setBallForm(prev => ({ ...prev, run_type: value as RunType }))}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {getRunTypeOptions().map(option => (
                                            <SelectItem key={option.value} value={option.value}>
                                                {option.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            <div>
                                <Label htmlFor="byes">Byes</Label>
                                <Input
                                    type="number"
                                    min="0"
                                    max="6"
                                    value={ballForm.byes || 0}
                                    onChange={(e) => setBallForm(prev => ({ ...prev, byes: parseInt(e.target.value) || 0 }))}
                                />
                            </div>

                            <div className="md:col-span-2">
                                <div className="flex items-center space-x-2">
                                    <input
                                        type="checkbox"
                                        id="is_wicket"
                                        checked={ballForm.is_wicket || false}
                                        onChange={(e) => setBallForm(prev => ({ ...prev, is_wicket: e.target.checked }))}
                                    />
                                    <Label htmlFor="is_wicket">Is Wicket</Label>
                                </div>
                            </div>

                            {ballForm.is_wicket && (
                                <div>
                                    <Label htmlFor="wicket_type">Wicket Type</Label>
                                    <Select
                                        value={ballForm.wicket_type}
                                        onValueChange={(value) => setBallForm(prev => ({ ...prev, wicket_type: value }))}
                                    >
                                        <SelectTrigger>
                                            <SelectValue />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {getWicketTypeOptions().map(option => (
                                                <SelectItem key={option.value} value={option.value}>
                                                    {option.label}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                </div>
                            )}
                        </div>

                        <div className="flex space-x-2 mt-6">
                            <Button
                                onClick={handleAddBall}
                                disabled={scoring}
                            >
                                {scoring ? 'Adding...' : 'Add Ball'}
                            </Button>
                            <Button
                                variant="outline"
                                onClick={() => setShowAddBallForm(false)}
                            >
                                Cancel
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}

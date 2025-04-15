import { useState } from 'react'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'

function WordQuizContent ({ volume, onVolumeChange, level, onLevelChange }) {
  const [currentQuiz, setCurrentQuiz] = useState(null)
  const [quizChoices, setQuizChoices] = useState([])

  const handleStartQuiz = level => {
    const numberOfChoices = 10
    const quizData = window.CreateQuiz(parseInt(level, 10), numberOfChoices)
    console.log(quizData)
    const choicesData = window.CreateQuizChoices()
    console.log(choicesData)
    setCurrentQuiz(quizData)
    setQuizChoices(choicesData)
  }

  return (
    <>
      <h1>単語クイズ</h1>
      <button onClick={() => handleStartQuiz(level)}>Start Quiz</button>

      <VolumeControl volume={volume} onVolumeChange={onVolumeChange} />
      <LevelControl level={level} onLevelChange={onLevelChange} />
    </>
  )
}
export default WordQuizContent

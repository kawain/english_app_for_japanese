import { useState, useEffect, useRef } from 'react'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'
import QuizChoices from './components/QuizChoices.jsx'
import { tts } from './utils/tts'
import { addExcludedWordId } from './Storage.jsx'

const numberOfChoices = 10

function WordQuizContent ({ volume, onVolumeChange, level, onLevelChange }) {
  const [currentQuiz, setCurrentQuiz] = useState(null)
  const [quizChoices, setQuizChoices] = useState([])
  const [isQuizStarted, setIsQuizStarted] = useState(false)
  const [selectedChoiceId, setSelectedChoiceId] = useState(null)
  const [answerResult, setAnswerResult] = useState(null)
  const [correctCount, setCorrectCount] = useState(0)
  const [totalQuestions, setTotalQuestions] = useState(0)
  const [isHighlighted, setIsHighlighted] = useState(false)
  const prevLevelRef = useRef(level)

  // クイズ中の単語読み上げとハイライト効果
  useEffect(() => {
    if (currentQuiz && isQuizStarted) {
      setIsHighlighted(true)
      const timer = setTimeout(() => {
        setIsHighlighted(false)
      }, 500)

      tts(currentQuiz.en, 'en-US', volume)
        .then(() => console.log(`読み上げ完了: ${currentQuiz.en}`))
        .catch(error => console.error('TTS error:', error))

      return () => clearTimeout(timer)
    }
  }, [currentQuiz, isQuizStarted, volume]) // volume も依存配列に追加

  // レベル変更時にクイズをリセットする useEffect
  useEffect(() => {
    // isQuizStarted が true で、かつレベルが実際に変更された場合にリセット
    if (isQuizStarted && prevLevelRef.current !== level) {
      console.log(
        `Level changed from ${prevLevelRef.current} to ${level}. Resetting quiz state.`
      )
      setIsQuizStarted(false) // スタート画面に戻す
      // 必要に応じて他の関連 state もリセット
      setCurrentQuiz(null)
      setQuizChoices([])
      setSelectedChoiceId(null)
      setAnswerResult(null)
      setCorrectCount(0)
      setTotalQuestions(0)
    }
    // 現在のレベルを ref に保存して次回の比較に使用
    prevLevelRef.current = level
  }, [level, isQuizStarted]) // level と isQuizStarted を依存関係に追加

  const fetchQuizData = currentLevel => {
    try {
      const quizData = window.CreateQuiz(
        parseInt(currentLevel, 10),
        numberOfChoices
      )
      console.log('Quiz Data:', quizData)
      const choicesData = window.CreateQuizChoices()
      console.log('Choices Data:', choicesData)

      if (!quizData || !choicesData) {
        console.error('クイズデータの取得に失敗しました。')
        return null
      }
      if (quizData.id === undefined) {
        console.error(
          'currentQuiz に id プロパティが含まれていません。',
          quizData
        )
        return null
      }
      if (choicesData.some(choice => choice.id === undefined)) {
        console.error(
          'quizChoices の一部要素に id プロパティが含まれていません。',
          choicesData
        )
        return null
      }
      return { quizData, choicesData }
    } catch (error) {
      console.error('Error fetching quiz data:', error)
      // エラー発生時にも null を返すなど、適切なエラーハンドリングを行う
      return null
    }
  }

  const handleStartQuiz = startLevel => {
    const data = fetchQuizData(startLevel)
    if (data) {
      setCurrentQuiz(data.quizData)
      setQuizChoices(data.choicesData)
      setIsQuizStarted(true)
      setSelectedChoiceId(null)
      setAnswerResult(null)
      // 新しいレベルで開始するので統計情報をリセット
      setCorrectCount(0)
      setTotalQuestions(0)
    } else {
      // データ取得失敗時の処理（例：エラーメッセージ表示）
      alert(
        'クイズデータの読み込みに失敗しました。ページを再読み込みしてください。'
      )
      setIsQuizStarted(false) // 開始失敗時はスタート画面のまま
    }
  }

  const handleChoiceChange = event => {
    if (answerResult !== null) return // 回答後は選択を変更できないようにする
    setSelectedChoiceId(event.target.value)
    console.log('Selected Choice ID:', event.target.value)
  }

  const handleAnswerSubmit = () => {
    if (!currentQuiz || currentQuiz.id === undefined) {
      console.error('正解のID(currentQuiz.id)が取得できません。')
      // ユーザーにエラーを通知することも検討
      alert('エラーが発生しました。問題データを取得できませんでした。')
      return
    }

    let result
    let isCorrect = false

    if (selectedChoiceId === null) {
      result = 'passed' // 未選択の場合はパス扱い
    } else if (String(selectedChoiceId) === String(currentQuiz.id)) {
      result = 'correct'
      isCorrect = true
    } else {
      result = 'incorrect'
    }

    setAnswerResult(result)
    setTotalQuestions(prev => prev + 1) // パスした場合も問題数にはカウント

    if (isCorrect) {
      setCorrectCount(prev => prev + 1)
      // ローカルストレージへの保存処理（エラーハンドリングを追加）
      try {
        addExcludedWordId(String(currentQuiz.id))
        window.AddStorage(currentQuiz.id) // この関数が存在し、適切に動作することが前提
      } catch (error) {
        console.error('Error saving to storage:', error)
        // ストレージ保存失敗時の処理（例：ユーザー通知）
      }
    }

    // 正誤に関わらず、回答後に例文などを読み上げる
    if (currentQuiz.en2) {
      tts(currentQuiz.en2, 'en-US', volume)
        .then(() => console.log(`読み上げ完了: ${currentQuiz.en2}`))
        .catch(error => console.error('TTS error (en2):', error))
    } else {
      // en2 がない場合でも、正解/不正解のフィードバックを音声で行うことも検討可能
      // 例: 正解なら "Correct!", 不正解なら "Incorrect." など
      const feedbackWord = isCorrect
        ? 'Correct'
        : result === 'passed'
        ? 'Passed'
        : 'Incorrect'
      tts(feedbackWord, 'en-US', volume).catch(error =>
        console.error('TTS error (feedback):', error)
      )
    }
  }

  const handleNextQuiz = () => {
    const data = fetchQuizData(level) // 現在のレベルで次の問題を取得
    if (data) {
      setCurrentQuiz(data.quizData)
      setQuizChoices(data.choicesData)
      setSelectedChoiceId(null) // 選択肢をリセット
      setAnswerResult(null) // 回答結果をリセット
    } else {
      console.error('次の問題の取得に失敗しました。')
      alert('次の問題の読み込みに失敗しました。')
      // 失敗した場合、クイズを終了（スタート画面に戻る）させるか、リトライを促すかなど検討
      setIsQuizStarted(false)
    }
  }

  const calculateAccuracy = () => {
    if (totalQuestions === 0) return 0
    return Math.round((correctCount / totalQuestions) * 100)
  }

  return (
    <>
      <div className='quiz-container'>
        {!isQuizStarted ? (
          // スタートボタン: クリックで現在のレベルでクイズを開始
          <button onClick={() => handleStartQuiz(level)}>
            単語クイズ・スタート (レベル {level})
          </button>
        ) : (
          // クイズ中の表示
          <>
            {currentQuiz && (
              <h2
                className={isHighlighted ? 'highlight' : ''}
                onClick={() => tts(currentQuiz.en, 'en-US', volume)}
                style={{ cursor: 'pointer' }}
                title='クリックして再読み上げ'
              >
                {currentQuiz.en}
              </h2>
            )}

            <div className='quiz-content'>
              {/* 選択肢表示 */}
              <QuizChoices
                quizChoices={quizChoices}
                selectedChoiceId={selectedChoiceId}
                onChoiceChange={handleChoiceChange}
                disabled={answerResult !== null}
              />
              <div>
                {answerResult === null ? (
                  // 「答える」ボタン
                  <button onClick={handleAnswerSubmit}>答える / パス</button>
                ) : (
                  <>
                    <div className='quiz-result'>
                      {answerResult === 'correct' && (
                        <p style={{ color: 'green', fontWeight: 'bold' }}>
                          正解です！
                        </p>
                      )}
                      {answerResult === 'incorrect' && (
                        <p style={{ color: 'red', fontWeight: 'bold' }}>
                          間違いです...
                        </p>
                      )}
                      {answerResult === 'passed' && (
                        <p style={{ color: 'orange', fontWeight: 'bold' }}>
                          パスしました
                        </p>
                      )}

                      <div className='quiz-details'>
                        <p>{currentQuiz.en}</p>
                        <p>{currentQuiz.jp}</p>
                        <p>{currentQuiz.en2}</p>
                        <p>{currentQuiz.jp2}</p>
                      </div>

                      <div className='quiz-stats'>
                        <p>
                          正解数: {correctCount} / {totalQuestions}
                        </p>
                        <p>正解率: {calculateAccuracy()}%</p>
                      </div>
                    </div>

                    <button onClick={handleNextQuiz}>次の問題を行う</button>
                  </>
                )}
              </div>
            </div>
          </>
        )}
      </div>
      <VolumeControl volume={volume} onVolumeChange={onVolumeChange} />
      <LevelControl level={level} onLevelChange={onLevelChange} />
    </>
  )
}
export default WordQuizContent

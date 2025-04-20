import { useState, useEffect, useCallback } from 'react'
import { useAppContext, speakText } from './App.jsx'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'
import { addExcludedWordId } from './Storage.jsx'

function ListeningContent () {
  const { selectedLevel, volume, isSoundEnabled } = useAppContext()
  const [progress, setProgress] = useState(0)
  const [times, setTimes] = useState(0)
  const [currentQuestion, setCurrentQuestion] = useState(null)
  const [en, setEn] = useState('')
  const [jp, setJp] = useState('【日本語訳】')
  const [en2, setEn2] = useState('【例文】')
  const [jp2, setJp2] = useState('【例文の日本語訳】')
  const [step, setStep] = useState(0)
  const [autoPlay, setAutoPlay] = useState(false)

  // スタート
  const handleStart = () => {
    const question = fetchQuestion()
    if (!question) {
      console.error('Failed to fetch question.')
      return
    }
    setCurrentQuestion(question)
    setEn(question.en)
    setTimes(1)
    setProgress(1)
    setStep(1)
    speakText(question.en, 'en-US', volume, isSoundEnabled)
  }

  const handleEnClick = () => {
    if (currentQuestion) {
      speakText(currentQuestion.en, 'en-US', volume, isSoundEnabled)
    }
  }

  const handleJpClick = () => {
    if (currentQuestion && currentQuestion.jp && currentQuestion.en2) {
      setJp(currentQuestion.jp)
      setStep(prev => prev + 1)
      speakText(currentQuestion.jp, 'ja-JP', volume, isSoundEnabled)
    }
  }

  const handleEn2Click = () => {
    if (currentQuestion && currentQuestion.en2) {
      setEn2(currentQuestion.en2)
      setStep(prev => prev + 1)
      speakText(currentQuestion.en2, 'en-US', volume, isSoundEnabled)
    }
  }

  const handleJp2Click = () => {
    if (currentQuestion && currentQuestion.jp2) {
      setJp2(currentQuestion.jp2)
      setStep(prev => prev + 1)
      speakText(currentQuestion.jp2, 'ja-JP', volume, isSoundEnabled)
    }
  }

  const handleNext = () => {
    const question = fetchQuestion()
    if (!question) {
      console.error('Failed to fetch question.')
      return
    }
    setCurrentQuestion(question)
    setEn(question.en)
    setJp('【日本語訳】')
    setEn2('【例文】')
    setJp2('【例文の日本語訳】')
    setTimes(prev => prev + 1)
    setStep(1)
    speakText(question.en, 'en-US', volume, isSoundEnabled)
  }

  // 問題を取得する関数 (useCallbackでメモ化)
  const fetchQuestion = useCallback(() => {
    try {
      const questionData = window.GetListeningQuestion(
        parseInt(selectedLevel, 10)
      )
      if (!questionData) {
        throw new Error('問題データを取得できませんでした。')
      }
      console.log('New question:', questionData)
      return questionData
    } catch (err) {
      console.error('Error fetching question:', err)
      return null
    }
  }, [selectedLevel]) // selectedLevelが変わったら関数を再生成

  // 自動シーケンスの処理
  useEffect(() => {
    if (progress === 0 || !currentQuestion || !autoPlay) return

    const getWaitTime = (text, lang) => {
      const baseTime = 2000
      const timePerChar = lang === 'ja-JP' ? 380 : 140
      return Math.max(baseTime, text.length * timePerChar)
    }

    const proceedToNextStep = () => {
      if (step === 1) {
        handleJpClick()
      } else if (step === 2) {
        handleEn2Click()
      } else if (step === 3) {
        handleJp2Click()
      } else if (step === 4) {
        handleNext()
      }
    }

    let timerId
    let waitTime

    // 現在のステップに応じたテキストを取得し、待機時間を計算
    if (step === 1) {
      waitTime = getWaitTime(currentQuestion.en, 'en-US')
    } else if (step === 2) {
      waitTime = getWaitTime(currentQuestion.jp, 'ja-JP')
    } else if (step === 3) {
      waitTime = getWaitTime(currentQuestion.en2, 'en-US')
    } else if (step === 4) {
      waitTime = getWaitTime(currentQuestion.jp2, 'ja-JP')
    }

    timerId = setTimeout(() => {
      proceedToNextStep()
    }, waitTime)

    return () => {
      clearTimeout(timerId)
    }
  }, [step, progress, autoPlay])

  // レベルを変更
  useEffect(() => {
    if (progress === 0) return
    handleStart()
  }, [selectedLevel])

  let content = null

  if (progress === 0) {
    content = <button onClick={handleStart}>リスニング開始</button>
  } else if (progress === 1) {
    content = (
      <>
        <div className='listening-content'>
          <div className='number-area'>{times}回目</div>
          <div
            className={step === 1 ? 'en-area highlight' : 'en-area'}
            onClick={handleEnClick}
            style={{ cursor: 'pointer' }}
          >
            {en}
          </div>
          <div
            className={step === 2 ? 'jp-area highlight' : 'jp-area'}
            onClick={handleJpClick}
            style={{ cursor: 'pointer' }}
          >
            {jp}
          </div>
          <div
            className={step === 3 ? 'en2-area highlight' : 'en2-area'}
            onClick={handleEn2Click}
            style={{ cursor: 'pointer' }}
          >
            {en2}
          </div>
          <div
            className={step === 4 ? 'jp2-area highlight' : 'jp2-area'}
            onClick={handleJp2Click}
            style={{ cursor: 'pointer' }}
          >
            {jp2}
          </div>
        </div>
        <div className='button-container'>
          <button onClick={handleNext} disabled={autoPlay}>
            次の問題へ
          </button>
          <button onClick={() => setAutoPlay(prev => !prev)}>
            {autoPlay ? '自動再生をオフ' : '自動再生をオン'}
          </button>
          <button
            onClick={() => {
              addExcludedWordId(String(currentQuestion.id))
              window.AddStorage(currentQuestion.id)
              alert('ストレージに追加しました')
            }}
          >
            ストレージに追加
          </button>
        </div>
      </>
    )
  }

  return (
    <>
      <div className='listening-container'>{content}</div>
      <VolumeControl />
      <LevelControl />
    </>
  )
}
export default ListeningContent

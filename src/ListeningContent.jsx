import { useState, useEffect, useCallback, use } from 'react'
import { useAppContext } from './App.jsx'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'
import { tts } from './utils/tts'
import { addExcludedWordId } from './Storage.jsx'

function ListeningContent () {
  const { selectedLevel, volume, isSoundEnabled } = useAppContext()
  const [times, setTimes] = useState(0)
  const [currentQuestion, setCurrentQuestion] = useState(null)
  // 0: 初期状態
  // 1: en表示/読み上げ
  // 2: jp表示/読み上げ
  // 3: en2表示/読み上げ
  // 4: jp2表示/読み上げ
  const [progress, setProgress] = useState(0)

  // スタート
  const handleStart = () => {
    fetchQuestion()
    setTimes(times + 1)
    setProgress(1)
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
      setCurrentQuestion(questionData)
    } catch (err) {
      console.error('Error fetching question:', err)
      setCurrentQuestion(null)
    }
  }, [selectedLevel]) // selectedLevelが変わったら関数を再生成

  useEffect(() => {
    if (progress === 0) return

    let timerId = null

    const handleProgress = async () => {
      try {
        if (isSoundEnabled) {
          if (progress === 1) {
            await tts(currentQuestion.en, 'en-US', volume, isSoundEnabled)
            setProgress(2)
          } else if (progress === 2) {
            await tts(currentQuestion.jp, 'ja-JP', volume, isSoundEnabled)
            setProgress(3)
          } else if (progress === 3) {
            await tts(currentQuestion.en2, 'en-US', volume, isSoundEnabled)
            setProgress(4)
          } else if (progress === 4) {
            await tts(currentQuestion.jp2, 'ja-JP', volume, isSoundEnabled)
            fetchQuestion()
            setTimes(times + 1)
            setProgress(1)
          }
        } else {
          timerId = setTimeout(() => {
            if (progress === 1) {
              setProgress(2)
            } else if (progress === 2) {
              setProgress(3)
            } else if (progress === 3) {
              setProgress(4)
            } else if (progress === 4) {
              fetchQuestion()
              setTimes(times + 1)
              setProgress(1)
            }
          }, 3000)
        }
      } catch (error) {
        console.error('TTS エラー:', error)
      }
    }

    handleProgress()

    return () => {
      if (timerId) {
        clearTimeout(timerId)
      }
    }
  }, [progress, currentQuestion, isSoundEnabled])

  // レベルを変えたらリセット
  useEffect(() => {
    if (progress === 0) return
    setTimes(0)
    setProgress(0)
  }, [selectedLevel])

  let content = null

  if (progress === 0) {
    content = <button onClick={handleStart}>リスニング開始</button>
  } else if (progress === 1) {
    content = (
      <>
        <div className='listening-content'>
          <div className='number-area'>{times}回目</div>
          <div className='en-area'>{currentQuestion.en}</div>
          <div className='jp-area' onClick={() => setProgress(2)}>
            ?
          </div>
          <div className='en2-area'>?</div>
          <div className='jp2-area'>?</div>
        </div>
      </>
    )
  } else if (progress === 2) {
    content = (
      <>
        <div className='listening-content'>
          <div className='number-area'>{times}回目</div>
          <div className='en-area'>{currentQuestion.en}</div>
          <div className='jp-area'>{currentQuestion.jp}</div>
          <div className='en2-area' onClick={() => setProgress(3)}>
            ?
          </div>
          <div className='jp2-area'>?</div>
        </div>
      </>
    )
  } else if (progress === 3) {
    content = (
      <>
        <div className='listening-content'>
          <div className='number-area'>{times}回目</div>
          <div className='en-area'>{currentQuestion.en}</div>
          <div className='jp-area'>{currentQuestion.jp}</div>
          <div className='en2-area'>{currentQuestion.en2}</div>
          <div className='jp2-area' onClick={() => setProgress(4)}>
            ?
          </div>
        </div>
      </>
    )
  } else if (progress === 4) {
    content = (
      <>
        <div className='listening-content'>
          <div className='number-area'>{times}回目</div>
          <div className='en-area'>{currentQuestion.en}</div>
          <div className='jp-area'>{currentQuestion.jp}</div>
          <div className='en2-area'>{currentQuestion.en2}</div>
          <div className='jp2-area'>{currentQuestion.jp2}</div>
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

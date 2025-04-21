import { useState, useEffect, useCallback } from 'react' // useMemo を追加
import { useAppContext } from './App.jsx'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'
import { addExcludedWordId } from './Storage.jsx'

function ListeningContent () {
  const { selectedLevel, isSoundEnabled, speak, syncOrAsyncMode } =
    useAppContext()
  const [progress, setProgress] = useState(0)
  const [times, setTimes] = useState(0)
  const [currentQuestion, setCurrentQuestion] = useState(null)
  const [en, setEn] = useState('')
  const [jp, setJp] = useState('【日本語訳】')
  const [en2, setEn2] = useState('【例文】')
  const [jp2, setJp2] = useState('【例文の日本語訳】')
  const [step, setStep] = useState(0)
  const [autoPlay, setAutoPlay] = useState(false)
  const [reviewArray, setReviewArray] = useState([])

  // 問題を取得する関数 (useCallbackでメモ化)
  const fetchQuestion = useCallback(() => {
    try {
      const questionData = window.GetListeningQuestion(
        parseInt(selectedLevel, 10)
      )
      if (!questionData) {
        throw new Error('問題データを取得できませんでした。')
      }
      setReviewArray(prev => [...prev, questionData])
      return questionData
    } catch (err) {
      console.error('Error fetching question:', err)
      return null
    }
  }, [selectedLevel]) // selectedLevelが変わったら関数を再生成

  const next = (startFlag = false) => {
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
    if (startFlag) {
      setProgress(1)
      setTimes(1)
    } else {
      setTimes(prev => prev + 1)
    }
    setStep(1)
  }

  const handleStart = () => {
    next(true)
  }

  const handleEnClick = () => {
    if (currentQuestion) {
      setEn(currentQuestion.en)
      speak(currentQuestion.en, 'en-US')
      setStep(1)
    }
  }

  const handleJpClick = () => {
    if (currentQuestion) {
      setJp(currentQuestion.jp)
      speak(currentQuestion.jp, 'ja-JP')
      setStep(2)
    }
  }

  const handleEn2Click = () => {
    if (currentQuestion) {
      setEn2(currentQuestion.en2)
      speak(currentQuestion.en2, 'en-US')
      setStep(3)
    }
  }

  const handleJp2Click = () => {
    if (currentQuestion) {
      setJp2(currentQuestion.jp2)
      speak(currentQuestion.jp2, 'ja-JP')
      setStep(4)
    }
  }

  const handleNext = () => {
    next()
  }

  useEffect(() => {
    if (
      progress === 0 ||
      !currentQuestion ||
      !autoPlay ||
      step === 0 ||
      !isSoundEnabled
    )
      return

    let isCancelled = false
    let timerId = null

    // 同期モード用の次のステップに進む関数
    const proceedToNextStepSync = () => {
      if (isCancelled) return
      // 次のステップのUI更新と読み上げ開始
      if (step === 1) {
        setJp(currentQuestion.jp)
        speak(currentQuestion.jp, 'ja-JP')
        setStep(2)
      } else if (step === 2) {
        setEn2(currentQuestion.en2)
        speak(currentQuestion.en2, 'en-US')
        setStep(3)
      } else if (step === 3 && currentQuestion.jp2) {
        setJp2(currentQuestion.jp2)
        speak(currentQuestion.jp2, 'ja-JP')
        setStep(4)
      } else if (step === 4) {
        next()
      }
    }

    if (syncOrAsyncMode === 'sync') {
      // 同期モード: setTimeoutで待機
      const getWaitTime = (text, lang) => {
        const baseTime = 2000
        const timePerChar = lang === 'ja-JP' ? 300 : 100
        // テキストがない場合のデフォルト待機時間
        if (!text) return baseTime
        return Math.max(baseTime, text.length * timePerChar)
      }

      let waitTime
      // 現在のステップの *読み上げが終わるまでの* 時間を見積もる
      if (step === 1 && currentQuestion.en)
        waitTime = getWaitTime(currentQuestion.en, 'en-US')
      else if (step === 2 && currentQuestion.jp)
        waitTime = getWaitTime(currentQuestion.jp, 'ja-JP')
      else if (step === 3 && currentQuestion.en2)
        waitTime = getWaitTime(currentQuestion.en2, 'en-US')
      else if (step === 4 && currentQuestion.jp2)
        waitTime = getWaitTime(currentQuestion.jp2, 'ja-JP')
      else return // 待機時間が計算できない場合は終了

      timerId = setTimeout(proceedToNextStepSync, waitTime)
    } else {
      // 非同期モード: speakの完了を待って次のステップへ
      const runAsyncSequence = async () => {
        if (isCancelled) return
        try {
          // このeffectはstepが変わった後に実行される
          // 現在のstepに応じた「次のアクション」（読み上げ＋state更新）を実行
          if (step === 1) {
            await speak(currentQuestion.en, 'en-US')
            await speak(currentQuestion.en, 'en-US')
            if (!isCancelled) {
              setJp(currentQuestion.jp)
              setStep(2)
            }
          } else if (step === 2) {
            await speak(currentQuestion.jp, 'ja-JP')
            if (!isCancelled) {
              setEn2(currentQuestion.en2)
              setStep(3)
            }
          } else if (step === 3) {
            await speak(currentQuestion.en2, 'en-US')
            if (!isCancelled) {
              setJp2(currentQuestion.jp2)
              setStep(4)
            }
          } else if (step === 4) {
            await speak(currentQuestion.jp2, 'ja-JP')
            await speak(currentQuestion.en2, 'en-US')
            await new Promise(resolve => setTimeout(resolve, 1000))
            if (!isCancelled) {
              next()
            }
          }
        } catch (error) {
          console.error('Error in async auto play sequence:', error)
          // エラーが発生した場合、自動再生を停止するなどの処理を追加可能
          if (!isCancelled) setAutoPlay(false)
        }
      }
      // 非同期シーケンスを開始
      runAsyncSequence()
    }

    // クリーンアップ関数
    return () => {
      isCancelled = true
      clearTimeout(timerId)
      // コンポーネントのアンマウント時や依存配列の値が変わる前に
      // 進行中の読み上げがあればキャンセルする
      if (window.speechSynthesis && window.speechSynthesis.speaking) {
        window.speechSynthesis.cancel()
      }
    }
  }, [
    step,
    progress,
    autoPlay,
    currentQuestion,
    syncOrAsyncMode,
    speak,
    isSoundEnabled
  ])

  // レベルを変更したときの処理
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
            onClick={() => {
              if (!autoPlay) {
                handleEnClick()
              }
            }}
            style={{ cursor: 'pointer' }}
          >
            {en}
          </div>
          <div
            className={step === 2 ? 'jp-area highlight' : 'jp-area'}
            onClick={() => {
              if (!autoPlay) {
                handleJpClick()
              }
            }}
            style={{ cursor: 'pointer' }}
          >
            {jp}
          </div>
          <div
            className={step === 3 ? 'en2-area highlight' : 'en2-area'}
            onClick={() => {
              if (!autoPlay) {
                handleEn2Click()
              }
            }}
            style={{ cursor: 'pointer' }}
          >
            {en2}
          </div>
          <div
            className={step === 4 ? 'jp2-area highlight' : 'jp2-area'}
            onClick={() => {
              if (!autoPlay) {
                handleJp2Click()
              }
            }}
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
              if (currentQuestion?.id != null) {
                addExcludedWordId(String(currentQuestion.id))
                window.AddStorage(currentQuestion.id)
                alert('ストレージに追加しました')
              } else {
                alert('現在の問題情報がありません。')
              }
            }}
            disabled={!currentQuestion}
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
      {reviewArray.length > 0 ? (
        <div className='review-container'>
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>English</th>
                <th>Japanese</th>
                <th>Example (EN/JA)</th>
                <th>Level</th>
                <th className='nowrap'>操作</th>
              </tr>
            </thead>
            <tbody>
              {reviewArray.map(item => (
                <tr key={item.id}>
                  <td className='center-text'>{item.id}</td>
                  <td>{item.en}</td>
                  <td>{item.jp}</td>
                  <td>
                    {item.en2}
                    <br />
                    {item.jp2}
                  </td>
                  <td className='center-text'>{item.level}</td>
                  <td className='center-text'>
                    <button
                      className='nowrap'
                      onClick={() => {
                        addExcludedWordId(String(currentQuestion.id))
                        window.AddStorage(currentQuestion.id)
                        const newArray = reviewArray.filter(
                          obj => obj.id !== item.id
                        )
                        setReviewArray(newArray)
                        alert('ストレージに追加しました')
                      }}
                    >
                      除外
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : null}

      <VolumeControl />
      <LevelControl />
    </>
  )
}
export default ListeningContent

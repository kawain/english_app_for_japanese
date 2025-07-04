import { useState, useEffect, useCallback } from 'react'
import { useAppContext } from './App.jsx'
import VolumeControl from './components/VolumeControl.jsx'

function ListeningContent2 () {
  const { isSoundEnabled, speak } = useAppContext()
  const [progress, setProgress] = useState(0)
  const [questionsList, setQuestionsList] = useState([]) // 取得した問題の配列全体を保持
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0) // 現在の問題のインデックス
  const [times, setTimes] = useState(0)
  const [currentQuestion, setCurrentQuestion] = useState(null)
  const [en2, setEn2] = useState('')
  const [jp2, setJp2] = useState('【日本語訳】')
  const [numberOfQuestionsToFetch, setNumberOfQuestionsToFetch] = useState(500) // 取得する問題数
  const [step, setStep] = useState(0)
  const [autoPlay, setAutoPlay] = useState(false)

  // WakeLock は画面起動ロック API のインターフェイスで、
  // アプリケーションが動作し続ける必要があるときに、
  // 端末の画面が暗くなったりロックされたりすることを防ぐためのものです。
  const [wakeLock, setWakeLock] = useState(null)
  const [isLocked, setIsLocked] = useState(false)

  // 問題を取得する関数 (useCallbackでメモ化)
  const fetchQuestion = useCallback(async count => {
    try {
      // GetListeningData は問題オブジェクトの配列を返す
      const questionsArray = await window.GetListeningData(count)
      // 配列が取得できたか、空でないかを確認
      if (!questionsArray || questionsArray.length === 0) {
        throw new Error(
          '問題データを取得できませんでした。または、該当する問題がありません。'
        )
      }
      return questionsArray
    } catch (err) {
      console.error('Error fetching questions:', err)
      return null
    }
  }, [])

  // next関数は、現在の問題リストとインデックスに基づいて問題を表示・進行させる
  // startFlag: true の場合、リストの最初の問題から開始 (またはリセット)
  // initialQuestions: startFlag が true の場合に参照する新しい問題リスト
  const next = useCallback(
    async (startFlag = false, initialQuestions = null) => {
      let questionToDisplay = null

      if (startFlag) {
        if (!initialQuestions || initialQuestions.length === 0) {
          console.error(
            'Failed to start: No questions provided or list is empty.'
          )
          // 問題がない場合はクリア
          setCurrentQuestion(null)
          return
        }
        questionToDisplay = initialQuestions[0]
        setCurrentQuestionIndex(0)
        setProgress(1)
        setTimes(1)
      } else {
        if (questionsList.length === 0) {
          console.error('Failed to get next question: Questions list is empty.')
          return
        }
        const newIndex =
          currentQuestionIndex + 1 >= questionsList.length
            ? 0 // リストの最後に達したら最初に戻る (ループ)
            : currentQuestionIndex + 1
        setCurrentQuestionIndex(newIndex)
        questionToDisplay = questionsList[newIndex]
        setTimes(prev => prev + 1)
      }

      setCurrentQuestion(questionToDisplay)
      setEn2(questionToDisplay.en2)
      setJp2('【日本語訳】')
      setStep(1)
    },
    [questionsList, currentQuestionIndex]
  )

  const handleStart = async () => {
    const newQuestions = await fetchQuestion(numberOfQuestionsToFetch)
    if (newQuestions && newQuestions.length > 0) {
      setQuestionsList(newQuestions)
      await next(true, newQuestions) // 新しい問題リストで開始
    } else {
      // 問題が取得できなかった場合の処理 (例: エラー表示)
      alert('開始できる問題がありません。')
      setProgress(1)
      setCurrentQuestion(null) // 現在の問題をクリア
    }
  }

  const handleEnClick = async () => {
    if (currentQuestion) {
      setStep(1)
      setEn2(currentQuestion.en2)
      await speak(currentQuestion.en2, 'en-US')
    }
  }

  const handleJpClick = async () => {
    if (currentQuestion) {
      setStep(2)
      setJp2(currentQuestion.jp2)
      await speak(currentQuestion.jp2, 'ja-JP')
    }
  }

  const handleNext = async () => {
    await next()
  }

  const handleAutoPlay = async () => {
    if (!isLocked) {
      try {
        const lock = await navigator.wakeLock.request('screen')
        setWakeLock(lock)
        setIsLocked(true)
        console.log('Screen Wake Lock acquired')
      } catch (err) {
        console.log(`${err.name}, ${err.message}`)
      }
    } else {
      if (wakeLock) {
        await wakeLock.release()
        setWakeLock(null)
        setIsLocked(false)
        console.log('Screen Wake Lock released')
      }
    }
    setAutoPlay(prev => !prev)
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

    // 非同期モード: speakの完了を待って次のステップへ
    const runAsyncSequence = async () => {
      try {
        // このeffectはstepが変わった後に実行される
        // 現在のstepに応じた「次のアクション」（読み上げ＋state更新）を実行
        if (step === 1) {
          await speak(currentQuestion.en2, 'en-US')
          await new Promise(resolve => setTimeout(resolve, 3000))
          setJp2(currentQuestion.jp2)
          setStep(2)
        } else if (step === 2) {
          await speak(currentQuestion.jp2, 'ja-JP')
          await speak(currentQuestion.en2, 'en-US')
          await new Promise(resolve => setTimeout(resolve, 3000))
          next()
        }
      } catch (error) {
        console.error('Error in async auto play sequence:', error)
      }
    }
    // 非同期シーケンスを開始
    runAsyncSequence()

    // クリーンアップ関数
    return () => {
      // 進行中の読み上げがあればキャンセルする
      if (window.speechSynthesis && window.speechSynthesis.speaking) {
        window.speechSynthesis.cancel()
      }
    }
  }, [step, progress, autoPlay, currentQuestion, speak, isSoundEnabled, next])

  let content = null

  if (progress === 0) {
    content = (
      <>
        <div className='start-container'>
          <button onClick={handleStart}>リスニング開始</button>
          <label htmlFor='numberOfQuestions'>使用する数: </label>
          <input
            type='number'
            id='numberOfQuestions'
            value={numberOfQuestionsToFetch}
            onChange={e =>
              setNumberOfQuestionsToFetch(parseInt(e.target.value, 10) || 1)
            }
            min='1'
            style={{ padding: '5px', width: '100px', textAlign: 'center' }}
          />
        </div>
      </>
    )
  } else if (progress === 1) {
    content = (
      <>
        <div className='listening-content'>
          <div className='number-area'>{times}回目</div>
          <div
            className={step === 1 ? 'en-area2 highlight' : 'en-area2'}
            onClick={() => {
              if (!autoPlay) {
                handleEnClick()
              }
            }}
            style={{ cursor: 'pointer' }}
          >
            {en2}
          </div>
          <div
            className={step === 2 ? 'jp-area highlight' : 'jp-area'}
            onClick={() => {
              if (!autoPlay) {
                handleJpClick()
              }
            }}
            style={{ cursor: 'pointer', marginBottom: '20px' }}
          >
            {jp2}
          </div>
        </div>
        <div className='button-container'>
          <button onClick={handleNext} disabled={autoPlay}>
            次の問題へ
          </button>
          <button onClick={handleAutoPlay} disabled={!isSoundEnabled}>
            {autoPlay ? '自動再生をオフ' : '自動再生をオン'}
          </button>
          <button
            disabled={autoPlay}
            onClick={async () => {
              if (currentQuestion?.id != null) {
                await window.AddStorage(currentQuestion.id)
                const newArray = questionsList.filter(
                  obj => obj.id !== currentQuestion.id
                )
                setQuestionsList(newArray)
                alert('ストレージに追加しました')
              } else {
                alert('現在の問題情報がありません。')
              }
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

      {questionsList.length > 0 ? (
        <div className='review-container'>
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>English</th>
                <th>Japanese</th>
                <th>Example (EN/JA)</th>
                <th className='nowrap'>操作</th>
              </tr>
            </thead>
            <tbody>
              {questionsList.map(item => (
                <tr key={item.id}>
                  <td className='center-text'>{item.id}</td>
                  <td
                    style={{ cursor: 'pointer' }}
                    role='button'
                    title='読み上げ'
                    onClick={async () => await speak(item.en, 'en-US')}
                  >
                    {item.en}
                  </td>
                  <td>{item.jp}</td>
                  <td
                    style={{ cursor: 'pointer' }}
                    role='button'
                    title='読み上げ'
                    onClick={async () => await speak(item.en2, 'en-US')}
                  >
                    {item.en2}
                    <br />
                    {item.jp2}
                  </td>
                  <td className='center-text'>
                    <button
                      className='nowrap'
                      onClick={async () => {
                        await window.AddStorage(item.id)
                        const newArray = questionsList.filter(
                          obj => obj.id !== item.id
                        )
                        setQuestionsList(newArray)
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
    </>
  )
}
export default ListeningContent2

import { useState, useEffect, useCallback } from 'react'
import VolumeControl from './components/VolumeControl.jsx'
import LevelControl from './components/LevelControl.jsx'
import { tts } from './utils/tts'
import { addExcludedWordId } from './Storage.jsx'

function ListeningContent ({
  level,
  onLevelChange,
  volume,
  onVolumeChange,
  isSoundEnabled,
  onToggleSound
}) {
  const [isListeningStarted, setIsListeningStarted] = useState(false)
  const [times, setTimes] = useState(0)
  const [currentQuestion, setCurrentQuestion] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState(null)
  // 0: 初期状態
  // 1: en表示/読み上げ
  // 2: jp表示/読み上げ
  // 3: en2表示/読み上げ
  // 4: jp2表示/読み上げ
  // 5: 完了
  const [displayStep, setDisplayStep] = useState(0)
  // 各要素の表示状態
  const [visibility, setVisibility] = useState({
    en: false,
    jp: false,
    en2: false,
    jp2: false
  })

  const handleStartListening = () => {
    setIsListeningStarted(true)
    setTimes(times + 1)
  }

  // 問題を取得する関数 (useCallbackでメモ化)
  const fetchQuestion = useCallback(() => {
    setIsLoading(true)
    setError(null)
    setDisplayStep(0) // ステップをリセット
    setVisibility({ en: false, jp: false, en2: false, jp2: false }) // 表示状態をリセット
    window.speechSynthesis.cancel() // 念のため既存の読み上げをキャンセル

    try {
      // window.GetListeningQuestion() は同期的にデータを返す想定
      const questionData = window.GetListeningQuestion(parseInt(level, 10))
      if (!questionData) {
        throw new Error('問題データを取得できませんでした。')
      }
      console.log('New question:', questionData)
      setCurrentQuestion(questionData)
      setDisplayStep(1) // 最初のステップへ移行
    } catch (err) {
      console.error('Error fetching question:', err)
      setError('問題の取得に失敗しました。')
      setCurrentQuestion(null)
    } finally {
      setIsLoading(false)
    }
  }, [level]) // levelが変わったら関数を再生成

  // コンポーネントマウント時とレベル変更時に最初の問題を取得
  useEffect(() => {
    fetchQuestion()
  }, [fetchQuestion]) // fetchQuestion (levelに依存) が変わったら実行

  // displayStep が変更されたら、表示と読み上げ処理を実行
  useEffect(() => {
    // ステップが 0 または 5以上(完了後) または問題がない場合は何もしない
    if (
      isListeningStarted === false ||
      displayStep === 0 ||
      displayStep > 4 ||
      !currentQuestion
    )
      return

    const processCurrentStep = async () => {
      try {
        let textToSpeak = ''
        let lang = ''
        let nextStep = displayStep + 1

        switch (displayStep) {
          case 1:
            if (currentQuestion.en) {
              setVisibility(prev => ({ ...prev, en: true }))
              textToSpeak = currentQuestion.en
              lang = 'en-US'
            }
            break
          case 2:
            if (currentQuestion.jp) {
              setVisibility(prev => ({ ...prev, jp: true }))
              textToSpeak = currentQuestion.jp
              lang = 'ja-JP'
            }
            break
          case 3:
            if (currentQuestion.en2) {
              setVisibility(prev => ({ ...prev, en2: true }))
              textToSpeak = currentQuestion.en2
              lang = 'en-US'
            }
            break
          case 4:
            if (currentQuestion.jp2) {
              setVisibility(prev => ({ ...prev, jp2: true }))
              textToSpeak = currentQuestion.jp2
              lang = 'ja-JP'
            }
            break
          default: // 完了状態にする
            // 想定外のステップ
            console.warn('Invalid display step:', displayStep)
            setDisplayStep(5)
            return
        }

        // テキストがあれば読み上げ、完了を待つ
        if (textToSpeak && lang) {
          await tts(textToSpeak, lang, volume, isSoundEnabled)
        } else {
          console.log(`Step ${displayStep}: No text to speak or invalid lang.`)
        }

        // 次のステップへ (エラーがなければ)
        setDisplayStep(nextStep)
      } catch (ttsError) {
        console.error(`TTS Error at step ${displayStep}:`, ttsError)
        // TTSエラーが発生した場合でも次のステップに進むか、処理を中断するか選択できます
        // ここでは、エラーが発生しても次のステップに進むようにしています
        setDisplayStep(prevStep => prevStep + 1)
      }
    }

    processCurrentStep()

    // クリーンアップ関数: displayStepが変わる前 or アンマウント時に読み上げをキャンセル
    return () => {
      window.speechSynthesis.cancel()
    }
  }, [displayStep, currentQuestion, volume, isSoundEnabled]) // 依存配列

  // 次の問題へ進む処理
  const handleNextQuestion = () => {
    fetchQuestion() // 新しい問題を取得（これによりuseEffectがトリガーされ、ステップがリセットされる）
  }

  // // 今は未使用後で使う
  // // 除外処理
  // const handleExclude = () => {
  //   if (currentQuestion && currentQuestion.id !== undefined) {
  //     try {
  //       addExcludedWordId(String(currentQuestion.id))
  //       window.AddStorage(currentQuestion.id) // WASM側にも通知
  //       console.log(`Excluded word ID: ${currentQuestion.id}`)
  //       // 除外したらすぐに次の問題へ
  //       fetchQuestion()
  //     } catch (error) {
  //       console.error('Error excluding word:', error)
  //       alert('単語の除外処理中にエラーが発生しました。')
  //     }
  //   } else {
  //     console.warn(
  //       'Cannot exclude: currentQuestion or currentQuestion.id is missing.'
  //     )
  //   }
  // }

  return (
    <>
      <div className='listening-container'>
        {!isListeningStarted ? (
          <button onClick={handleStartListening}>ヒアリング開始</button>
        ) : (
          <>
            <div className='listening-content'>
              <div className='number-area'>{times}回目</div>
              {visibility.en && (
                <div className='en-area'>{currentQuestion.en}</div>
              )}
              {visibility.jp && (
                <div className='jp-area'>{currentQuestion.jp}</div>
              )}
              {visibility.en2 && (
                <div className='en2-area'>{currentQuestion.en2}</div>
              )}
              {visibility.jp2 && (
                <div className='jp2-area'>{currentQuestion.jp2}</div>
              )}
            </div>
          </>
        )}
      </div>
      <VolumeControl
        volume={volume}
        onVolumeChange={onVolumeChange}
        isSoundEnabled={isSoundEnabled}
        onToggleSound={onToggleSound}
      />
      <LevelControl level={level} onLevelChange={onLevelChange} />
    </>
  )
}
export default ListeningContent

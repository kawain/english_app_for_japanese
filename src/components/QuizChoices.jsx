// props に selectedChoiceId と onChoiceChange を追加
function QuizChoices({ quizChoices, selectedChoiceId, onChoiceChange }) {
    // quizChoices が配列でない場合や空の場合のガード処理を追加
    if (!Array.isArray(quizChoices) || quizChoices.length === 0) {
      return <div>選択肢を読み込み中...</div>; // または他の適切な表示
    }
  
    return (
      <div className='quiz-choices'>
        {quizChoices.map((choice, index) => {
          // choice オブジェクトに必要なプロパティ (id, jp) があるか確認
          if (choice === null || choice.id === undefined || choice.jp === undefined) {
            console.warn(`選択肢 ${index} のデータが不正です:`, choice);
            return null; // 不正なデータはスキップ
          }
          const inputId = `choice-${choice.id}`; // idが重複しないようにプレフィックスを追加
          return (
            <div key={index} className='choice-item'>
              <input
                type='radio'
                id={inputId} // id属性値を設定
                name='quizChoice'
                value={choice.id}
                // 選択されているIDと一致するかどうかで checked を制御
                checked={String(selectedChoiceId) === String(choice.id)}
                // 変更時に onChoiceChange を呼び出す
                onChange={onChoiceChange}
              />
              {/* label の htmlFor も input の id と一致させる */}
              <label htmlFor={inputId}>{choice.jp}</label>
            </div>
          );
        })}
      </div>
    );
  }
  
  export default QuizChoices;
  
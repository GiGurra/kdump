package com.github.gigurra.kdump.internal.util

import scala.concurrent.duration.{Duration, FiniteDuration}
import scala.concurrent.{Await, ExecutionContext, Future, blocking}

object async {

  def apply[T](expr: => T)(using ec: ExecutionContext = ExecutionContext.global): asyncOp[T] =
    asyncOp(Future(blocking(expr))(ec))

  def await[T](f: asyncOp[T])(using ec: ExecutionContext = ExecutionContext.global): T =
    Await.result(f.op, Duration.Inf)

  case class asyncOp[T](op: Future[T]) {
    def join(using ec: ExecutionContext = ExecutionContext.global): T = this: T
  }

  given asyncOp2Val[T](using ec: ExecutionContext = ExecutionContext.global): Conversion[asyncOp[T], T] = x => await(x)
}
